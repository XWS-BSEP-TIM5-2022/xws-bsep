package api

import (
	"fmt"
	"log"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	authService       *AuthService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(authService *AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		authService:       authService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	var authRoles []domain.Role
	auth := &domain.Authentication{
		Id:       command.User.Id,
		Username: command.User.Username,
		Password: command.User.Password,
		Roles:    &authRoles,
	}
	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.CreateAuth:
		err := handler.authService.Register(*auth, command.User.Role, command.User.Email)
		fmt.Println(err)
		if err != nil {
			fmt.Println("Auth credentials are not saved, err:", err)
			reply.Type = events.AuthNotCreated
		} else {
			fmt.Println("Sve ok")
			reply.Type = events.AuthCreated
		}
	case events.RollbackAuth:
		err := handler.authService.Delete(auth.Id)
		if err != nil {
			log.Println("Auth credentials are not deleted, err: ", err)
		} else {
			reply.Type = events.AuthRolledBack
		}
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
