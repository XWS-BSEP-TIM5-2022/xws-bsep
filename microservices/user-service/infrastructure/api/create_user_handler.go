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
	reply := events.CreateUserReply{User: command.User}
	switch command.Type {
	case events.CreateUser:
		err := handler.userService.CheckEmailCriteria(command.User.Email)
		if err != nil {
			fmt.Println(err.Error())
			reply.Type = events.UserNotCreated
		} else {
			user := mapCommandUserToDomainUser(command)
			newUser, err := handler.userService.Insert(user)
			if err != nil {
				reply.Type = events.UserDeleted
				fmt.Println("User is not saved, err: ", err)
			} else {
				command.User.Id = newUser.Id.Hex()
				reply.Type = events.UserCreated
				reply.User = command.User
			}
		}
	case events.ApproveUser:
		objID, err := primitive.ObjectIDFromHex(command.User.Id)
		if err != nil {
			log.Println("User is not approved, err: ", err)
		} else {
			User := &domain.User{
				Id: objID,
			}
			fmt.Println("User is approved! ", User)
			reply.Type = events.UserApproved
			return
		}
	case events.DeleteUser:
		if command.User.Id == "" {
			log.Println("Ende saga ")
			return
		}
		objID, err := primitive.ObjectIDFromHex(command.User.Id)
		if err != nil {
			log.Println("User is not deleted $$$$, err: ", err)
		} else {
			User := &domain.User{
				Id:    objID,
				Email: command.User.Email,
			}
			err = handler.userService.Delete(User)
			if err != nil {
				log.Println("User is not deleted by id: ", err)
			} else {
				reply.Type = events.UserDeleted
				return
			}
		}
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
