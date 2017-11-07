package gosiris

import (
	"fmt"
	"log"
	"time"
)

type ActorRef struct {
	name        string
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

type ActorRefInterface interface {
	Tell(string, interface{}, ActorRefInterface) error
	Repeat(string, time.Duration, interface{}, ActorRefInterface) (chan struct{}, error)
	AskForClose(ActorRefInterface)
	LogInfo(string, ...interface{})
	LogError(string, ...interface{})
	Become(string, func(Message)) error
	Unbecome(string) error
	Name() string
	Forward(Message, ...string)
}

func newActorRef(name string) ActorRefInterface {
	ref := ActorRef{}
	ref.infoLogger, ref.errorLogger =
		NewActorLogger(name)
	ref.name = name
	return ref
}

func (ref ActorRef) LogInfo(format string, a ...interface{}) {
	logZipkinEvent(fmt.Sprintf(format, a...))
	ref.infoLogger.Printf(format, a...)
}

func (ref ActorRef) LogError(format string, a ...interface{}) {
	logZipkinEvent(fmt.Sprintf(format, a...))
	ref.errorLogger.Printf(format, a...)
}

func (ref ActorRef) Tell(messageType string, data interface{}, sender ActorRefInterface) error {
	actor, err := ActorSystem().actor(ref.name)

	if err != nil {
		ErrorLogger.Printf("Failed to send from %v to %v: %v", sender.Name(), ref.name, err)
		return err
	}

	dispatch(actor.actor.getDataChan(), messageType, data, &ref, sender, actor.options)

	return nil
}

func (ref ActorRef) Repeat(messageType string, d time.Duration, data interface{}, sender ActorRefInterface) (chan struct{}, error) {
	actor, err := ActorSystem().actor(ref.name)

	if err != nil {
		ErrorLogger.Printf("Failed to send from %v to %v: %v", sender.Name(), ref.name, err)
		return nil, err
	}

	t := time.NewTicker(d)
	stop := make(chan struct{})

	go
		func(t *time.Ticker, stop chan struct{}) {
			for {
				select {
				case <-t.C:
					dispatch(actor.actor.getDataChan(), messageType, data, &ref, sender, actor.options)
				case <-stop:
					t.Stop()
					close(stop)
					return
				}
			}
		}(t, stop)

	return stop, nil
}

func (ref ActorRef) AskForClose(sender ActorRefInterface) {
	InfoLogger.Printf("Asking to close %v", ref.name)

	actor, err := ActorSystem().actor(ref.name)

	if err != nil {
		InfoLogger.Printf("Actor %v already closed", ref.name)
		return
	}

	go dispatch(actor.actor.getDataChan(), GosirisMsgPoisonPill, nil, &ref, sender, actor.options)
}

func (ref ActorRef) Become(messageType string, f func(Message)) error {
	actor, err := ActorSystem().actor(ref.name)
	if err != nil {
		return fmt.Errorf("actor implementation %v not found", messageType)
	}

	if actor.actor.reactions() == nil {
		return fmt.Errorf("react for %v not yet implemented", messageType)
	}

	v, exists := actor.actor.reactions()[messageType]

	if !exists {
		return fmt.Errorf("react for %v not yet implemented", messageType)
	}

	actor.actor.unbecomeHistory()[messageType] = v
	actor.actor.reactions()[messageType] = f

	return nil
}

func (ref ActorRef) Unbecome(messageType string) error {
	actor, err := ActorSystem().actor(ref.name)
	if err != nil {
		return fmt.Errorf("actor implementation %v not found", messageType)
	}

	if actor.actor.reactions() == nil {
		return fmt.Errorf("become for %v not yet implemented", messageType)
	}

	v, exists := actor.actor.unbecomeHistory()[messageType]

	if !exists {
		return fmt.Errorf("unbecome for %v not yet implemented", messageType)
	}

	actor.actor.reactions()[messageType] = v
	delete(actor.actor.unbecomeHistory(), messageType)

	return nil
}

func (ref ActorRef) Name() string {
	return ref.name
}

func (ref ActorRef) Forward(message Message, destinations ...string) {
	for _, v := range destinations {
		actorRef, err := ActorSystem().ActorOf(v)
		if err != nil {
			fmt.Errorf("actor %v is not part of the actor system", v)
		}
		actorRef.Tell(message.MessageType, message.Data, message.Sender)
	}
}
