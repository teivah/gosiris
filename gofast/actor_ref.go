package gofast

import (
	"fmt"
	"log"
	"Gofast/gofast/util"
)

type ActorRef struct {
	name        string
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

type ActorRefInterface interface {
	Send(string, interface{}, ActorRefInterface) error
	AskForClose(ActorRefInterface)
	LogInfo(string, ...interface{})
	LogError(string, ...interface{})
	Become(string, func(Message)) error
	Unbecome(string) error
	Name() string
}

func newActorRef(name string) ActorRefInterface {
	ref := ActorRef{}
	ref.infoLogger, ref.errorLogger =
		util.NewActorLogger(name)
	ref.name = name
	return ref
}

func (ref ActorRef) LogInfo(format string, a ...interface{}) {
	ref.infoLogger.Printf(format, a...)
}

func (ref ActorRef) LogError(format string, a ...interface{}) {
	ref.errorLogger.Printf(format, a...)
}

func (ref ActorRef) Send(messageType string, data interface{}, sender ActorRefInterface) error {
	actor, err := ActorSystem().actor(ref)

	if err != nil {
		return err
	}

	dispatch(actor.Mailbox(), messageType, data, &ref, sender)

	return nil
}

func (ref ActorRef) AskForClose(sender ActorRefInterface) {
	util.LogError("Asking to close %v", ref.name)

	actor, err := ActorSystem().actor(ref)

	if err != nil {
		util.LogInfo("Actor %v already closed", ref.name)
		return
	}

	dispatch(actor.Mailbox(), poisonPill, nil, &ref, sender)
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

func (ref ActorRef) Name() string {
	return ref.name
}
