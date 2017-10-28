package gofast

import (
	"testing"
	"time"
)

type ParentActor struct {
	Actor
}

type ChildActor struct {
	Actor
	hello map[string]bool
}

func TestBasic(t *testing.T) {
	//Create a simple parent actor
	parentActor := Actor{}
	//Register the parent actor
	ActorSystem().RegisterActor("parentActor", &parentActor, Root())

	//Create a simple child actor
	childActor := Actor{}
	//Register the reactions to event types (here a reaction to message)
	childActor.React("message", func(message Message) {
		childActor.Printf("Received %v", message.Data)
	})
	//Register the child actor
	ActorSystem().RegisterActor("childActor", &childActor, &parentActor)

	//Retrieve the parent and child actor reference
	parentActorRef, _ := ActorSystem().Actor("parentActor")
	childActorRef, _ := ActorSystem().Actor("childActor")

	//Send a message from the parent to the child actor
	childActorRef.Tell("message", "Hi! How are you?", parentActorRef)
}

func TestStatefulness(t *testing.T) {
	childActor := ChildActor{}
	childActor.hello = make(map[string]bool)
	parentActor := ParentActor{}

	f := func(message Message) {
		parentActor.Printf("Receive response %v\n", message.Data)
	}

	parentActor.React("helloback", f).React("error", f).React("help", f)
	ActorSystem().RegisterActor("parent", &parentActor, Root())

	childActor.React("hello", func(message Message) {
		childActor.Printf("Receive request %v\n", message.Data)

		name := message.Data.(string)

		if _, ok := childActor.hello[name]; ok {
			message.Sender.Tell("error", "I already know you!", message.Self)
			childActor.Parent().Tell("help", "Daddy help me!", message.Self)
			childActor.Close()
		} else {
			childActor.hello[message.Data.(string)] = true
			message.Sender.Tell("helloback", "hello "+name+"!", message.Self)
		}
	})
	ActorSystem().RegisterActor("child", &childActor, &parentActor)

	childActorRef, _ := ActorSystem().Actor("child")
	parentActorRef, _ := ActorSystem().Actor("parent")

	childActorRef.Tell("hello", "teivah", parentActorRef)
	childActorRef.Tell("hello", "teivah", parentActorRef)

	time.Sleep(1000)
}
