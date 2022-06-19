package api

import (
	"fmt"
	"log"

	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserCommandHandler struct {
	userService       *application.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(userService *application.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	u := &CreateUserCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := u.commandSubscriber.Subscribe(u.handle)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	fmt.Println(" @@@@@@@@@@ hendleeer! ID: ", command.User.Id)
	reply := events.CreateUserReply{User: command.User}
	fmt.Println(command.Type)

	switch command.Type {
	case events.CreateUser:
		err := handler.userService.CheckEmailCriteria(command.User.Email)
		if err != nil {
			fmt.Println(err.Error())
			reply.Type = events.UserNotCreated
			return
		}

		user := mapCommandUserToDomainUser(command)
		newUser, err := handler.userService.Insert(user)
		if err != nil {
			reply.Type = events.UserNotCreated
			fmt.Println("GRESKA PRILIKOM INSERTA: ", err)
			return
		}

		command.User.Id = newUser.Id.Hex()
		log.Println("#### #### *************** ID", command.User.Id)
		reply.Type = events.UserCreated
		reply.User = command.User
		fmt.Println("Reply type: User created ", reply.Type, reply.User)

	case events.ApproveUser:
		fmt.Println("COMMAND USER ID", command.User.Id)
		objID, err := primitive.ObjectIDFromHex(command.User.Id)
		if err != nil {
			panic(err)
		}
		User := &domain.User{
			Id: objID,
		}
		err = handler.userService.Approve(User)
		if err != nil {
			// reply.Type = events.UserApproved
			return
		}
		reply.Type = events.UserApproved

	case events.DeleteUser:
		objID, err := primitive.ObjectIDFromHex(command.User.Id)
		if err != nil {
			panic(err)
		}
		User := &domain.User{
			Id: objID,
		}
		err = handler.userService.Delete(User)
		if err != nil {
			return
		}
		reply.Type = events.UserDeleted

	case events.RollbackUser:

		fmt.Println("TODO SD: ROLLBACK USER")
		reply.Type = events.UserRolledBack

	default:
		reply.Type = events.UnknownReply
	}

	log.Println("ODGOVOR RBR: ", reply.Type)
	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
