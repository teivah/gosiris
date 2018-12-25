package gosiris

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"log"
	"time"
)

type ActorRef struct {
	name        string
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

type ActorRefInterface interface {
	Tell(Context, string, interface{}, ActorRefInterface) error
	Repeat(string, time.Duration, interface{}, ActorRefInterface) (chan struct{}, error)
	AskForClose(ActorRefInterface)
	LogInfo(Context, string, ...interface{})
	LogError(Context, string, ...interface{})
	Become(string, func(Context)) error
	Unbecome(string) error
	Name() string
	Forward(Context, ...string)
	Ask(Context, string, interface{}, time.Duration) (interface{}, error)
}

func newActorRef(name string) ActorRefInterface {
	ref := ActorRef{}
	ref.infoLogger, ref.errorLogger =
		NewActorLogger(name)
	ref.name = name
	return ref
}

func (ref ActorRef) LogInfo(context Context, format string, a ...interface{}) {
	ref.infoLogger.Printf(format, a...)
	if zipkinSystemInitialized {
		logZipkinMessage(context.span, fmt.Sprintf(format, a...))
	}
}

func (ref ActorRef) LogError(context Context, format string, a ...interface{}) {
	ref.errorLogger.Printf(format, a...)
	if zipkinSystemInitialized {
		logZipkinMessage(context.span, fmt.Sprintf(format, a...))
	}
}

func (ref ActorRef) Tell(context Context, messageType string, data interface{}, sender ActorRefInterface) error {
	actor, err := ActorSystem().actor(ref.name)

	if err != nil {
		ErrorLogger.Printf("Failed to send from %v to %v: %v", sender.Name(), ref.name, err)
		return err
	}

	var span opentracing.Span = nil
	if context.span != nil {
		span = context.span
	} else if zipkinSystemInitialized && messageType != GosirisMsgChildClosed {
		span = startZipkinSpan(sender.Name(), messageType)
	}

	dispatch(actor.actor.getDataChan(), messageType, data, &ref, sender, actor.options, span)

	if span != nil {
		stopZipkinSpan(span)
	}

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

	go func(t *time.Ticker, stop chan struct{}) {
		for {
			select {
			case <-t.C:
				dispatch(actor.actor.getDataChan(), messageType, data, &ref, sender, actor.options, nil)
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

	go dispatch(actor.actor.getDataChan(), GosirisMsgPoisonPill, nil, &ref, sender, actor.options, nil)
}

func (ref ActorRef) Become(messageType string, f func(Context)) error {
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

func (ref ActorRef) Forward(context Context, destinations ...string) {
	for _, v := range destinations {
		actorRef, err := ActorSystem().ActorOf(v)
		if err != nil {
			ErrorLogger.Printf("actor %v is not part of the actor system", v)
		}
		actorRef.Tell(context, context.MessageType, context.Data, context.Sender)
	}
}

func (ref ActorRef) Ask(context Context, messageType string, data interface{}, timeout time.Duration) (interface{}, error) {
	temp := Actor{name: "temp_" + uuid.New().String()}
	defer temp.Close()

	ch := make(chan interface{}, 1)
	temp.React(messageType, func(ctx Context) {
		ch <- ctx.Data
	})

	err := ActorSystem().RegisterActor(temp.name, &temp, nil)
	if err != nil {
		ErrorLogger.Printf("Failed to ask to %v: %v", ref.name, err)
		return nil, err
	}

	tempRef, _ := ActorSystem().ActorOf(temp.name)
	err = ref.Tell(EmptyContext, messageType, data, tempRef)
	if err != nil {
		ErrorLogger.Printf("Failed to ask to %v: %v", ref.name, err)
		return nil, err
	}

	select {
	case reply := <-ch:
		return reply, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout while waiting for the answer")
	}
}
