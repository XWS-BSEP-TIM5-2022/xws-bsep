package api

import (
	"fmt"

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
	id, err := primitive.ObjectIDFromHex(command.User.Id)
	if err != nil {
		return
	}
	User := &domain.User{
		Id: id,
	}
	reply := events.CreateUserReply{User: command.User}

	fmt.Println(command.Type)
	switch command.Type {
	case events.CreateUser:
		fmt.Println("create user -> user created ")

		err := handler.userService.CheckEmailCriteria(command.User.Email)
		if err != nil {
			fmt.Println(err.Error())
			reply.Type = events.UserNotCreated
			return
		}
		// todo: map command to user
		user := domain.User{
			// Id: User.Id,
			Email: command.User.Email,
		}
		newUser, err := handler.userService.Insert(&user)
		if err != nil {
			reply.Type = events.UserNotCreated
			return
		}

		command.User.Id = newUser.Id.Hex()
		reply.Type = events.UserCreated
	case events.ApproveUser:
		err := handler.userService.Approve(User)
		if err != nil {
			// reply.Type = events.UserApproved
			return
		}
		reply.Type = events.UserApproved
	case events.DeleteUser:
		err := handler.userService.Delete(User)
		if err != nil {
			return
		}
		reply.Type = events.UserDeleted
	case events.RollbackUser:
		fmt.Println("TODO SD: ROLLBACK USER")
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
