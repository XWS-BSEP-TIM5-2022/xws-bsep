package api

import (
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	authService       *application.AuthService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(authService *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
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
	// TODO SD: logika
	var authRoles []domain.Role
	// for _, authRole := range command.User.Role {
	// 	roles, err := handler.authService.store.FindRoleByName(authRole)
	// 	if err != nil {
	// 		fmt.Println("Error finding role by name")
	// 		return nil, err
	// 	}
	// 	authRoles = append(authRoles, *roles...)
	// }
	auth := &domain.Authentication{
		Id:       command.User.Id,
		Username: command.User.Username,
		Password: command.User.Password,
		Roles:    &authRoles, // TODO SD: prazna lista rola
	}

	fmt.Println(" ************* " + command.User.Id + " ************* ")

	reply := events.CreateUserReply{User: command.User}

	fmt.Println(command.Type)
	switch command.Type {
	case events.CreateAuth:
		err := handler.authService.Register(*auth, command.User.Role, command.User.Email)
		fmt.Println("greska registracija saga: ", err)
		if err != nil {
			reply.Type = events.AuthNotCreated
			return
		}
		reply.Type = events.AuthCreated
	case events.DeleteUser:
		// TODO SD: ispraviti
		// err := handler.authService.Cancel(auth)
		// if err != nil {
		// 	return
		// }
		reply.Type = events.UserDeleted
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
