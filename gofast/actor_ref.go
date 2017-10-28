package gofast

import "fmt"

type ActorRef struct {
	name string
}

type ActorRefInterface interface {
	Send(string, interface{}, ActorRefInterface) error
	AskForClose(ActorRefInterface)
	Printf(string, ...interface{}) (int, error)
	Become(string, func(Message)) error
	Unbecome(string) error
}

func (ref ActorRef) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("["+ref.name+"] "+format, a...)
}

func (ref ActorRef) Send(messageType string, data interface{}, sender ActorRefInterface) error {
	actor, _ := ActorSystem().actor(ref)

	dispatch(actor.Mailbox(), messageType, data, &ref, sender)

	return nil
}

func (ref ActorRef) AskForClose(sender ActorRefInterface) {
	actor, err := ActorSystem().actor(ref)

	if err != nil {
		return
	}

	dispatch(actor.Mailbox(), poisonPill.messageType, poisonPill.Data, &ref, sender)
}

func (ref ActorRef) Become(messageType string, f func(Message)) error {
	actor, err := ActorSystem().actor(ref)
	if err != nil {
		return fmt.Errorf("actor implementation %v not found", messageType)
	}

	if actor.reactions() == nil {
		return fmt.Errorf("react for %v not yet implemented", messageType)
	}

	v, exists := actor.reactions()[messageType]

	if !exists {
		return fmt.Errorf("react for %v not yet implemented", messageType)
	}

	actor.Unbecome()[messageType] = v
	actor.reactions()[messageType] = f

	return nil
}

func (ref ActorRef) Unbecome(messageType string) error {
	actor, err := ActorSystem().actor(ref)
	if err != nil {
		return fmt.Errorf("actor implementation %v not found", messageType)
	}

	if actor.reactions() == nil {
		return fmt.Errorf("become for %v not yet implemented", messageType)
	}

	v, exists := actor.Unbecome()[messageType]

	if !exists {
		return fmt.Errorf("unbecome for %v not yet implemented", messageType)
	}

	actor.reactions()[messageType] = v
	delete(actor.Unbecome(), messageType)

	return nil
}
