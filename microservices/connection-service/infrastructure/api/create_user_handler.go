package api

import (
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
)

type CreateUserCommandHandler struct {
	userService *application.UserService
	// replyPublisher    saga.Publisher
	// commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(userService *application.UserService /*, publisher saga.Publisher, subscriber saga.Subscriber*/) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		userService: userService,
		// replyPublisher:    publisher,
		// commandSubscriber: subscriber,
	}
	fmt.Print(o)
	return o, nil
	// err := o.commandSubscriber.Subscribe(o.handle)
	// if err != nil {
	// 	return nil, err
	// }
	// return o, nil
}