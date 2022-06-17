package application

import (
	"fmt"

	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
)

type CreateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserOrchestrator, error) {
	o := &CreateUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *CreateUserOrchestrator) Start(event *events.CreateUserCommand) error {
	return o.commandPublisher.Publish(event)
}

func (o *CreateUserOrchestrator) handle(reply *events.CreateUserReply) {
	command := events.CreateUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateUserOrchestrator) nextCommandType(reply events.CreateUserReplyType) events.CreateUserCommandType {
	// fmt.Println(" *** Odgovor: ", reply)
	switch reply {
	case events.UserCreated:
		fmt.Println(" #UserCreated ", events.UserCreated)
		return events.CreateAuth
	case events.UserNotCreated:
		fmt.Println(" #DeleteUser ", events.DeleteUser)
		return events.DeleteUser
	case events.UserRolledBack:
		fmt.Println(" #UserRolledBack ", events.UserRolledBack)
		return events.DeleteUser
	case events.AuthNotCreated:
		fmt.Println(" # ROLLBACK! ", events.RollbackUser)
		return events.RollbackUser
	case events.AuthCreated:
		fmt.Println(" #ApproveUser  ", events.ApproveUser)
		return events.ApproveUser
	default:
		fmt.Println(" # UnknownCommand")
		return events.UnknownCommand
	}
}
