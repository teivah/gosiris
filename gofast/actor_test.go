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

func init() {
	InitRemoteActorSystem("http://192.168.99.100:2379")
}

func TestBasic(t *testing.T) {
	//Create a simple parent actor
	parentActor := Actor{}
	defer parentActor.Close()

	//Register the parent actor
	ActorSystem().RegisterActor("parentActor", &parentActor)

	//Create a simple child actor
	childActor := Actor{}
	defer childActor.Close()

	//Register the reactions to event types (here a reaction to message)
	childActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})

	//Register the child actor
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor)

	//Retrieve the parent and child actor reference
	parentActorRef, _ := ActorSystem().Actor("parentActor")
	childActorRef, _ := ActorSystem().Actor("childActor")

	//Send a message from the parent to the child actor
	childActorRef.Send("message", "Hi! How are you?", parentActorRef)
}

func TestStatefulness(t *testing.T) {
	childActor := ChildActor{}
	childActor.hello = make(map[string]bool)

	parentActor := ParentActor{}
	defer parentActor.Close()

	f := func(message Message) {
		message.Self.LogInfo("Receive response %v\n", message.Data)
	}

	parentActor.React("helloback", f).React("error", f).React("help", f)
	ActorSystem().RegisterActor("parent", &parentActor)

	childActor.React("hello", func(message Message) {
		message.Self.LogInfo("Receive request %v\n", message.Data)

		name := message.Data.(string)

		if _, ok := childActor.hello[name]; ok {
			message.Sender.Send("error", "I already know you!", message.Self)
			childActor.Parent().Send("help", "Daddy help me!", message.Self)
			childActor.Close()
		} else {
			childActor.hello[message.Data.(string)] = true
			message.Sender.Send("helloback", "hello "+name+"!", message.Self)
		}
	})
	ActorSystem().SpawnActor(&parentActor, "child", &childActor)

	childActorRef, _ := ActorSystem().Actor("child")
	parentActorRef, _ := ActorSystem().Actor("parent")

	childActorRef.Send("hello", "teivah", parentActorRef)
	childActorRef.Send("hello", "teivah", parentActorRef)

	time.Sleep(500 * time.Millisecond)
	childActorRef.Send("hello", "teivah2", parentActorRef)
}

func TestClose(t *testing.T) {
	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor)

	childActor := Actor{}
	childActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor)

	parentActorRef, _ := ActorSystem().Actor("parentActor")
	childActorRef, _ := ActorSystem().Actor("childActor")

	childActorRef.AskForClose(parentActorRef)

	time.Sleep(500 * time.Millisecond)
	childActorRef.Send("message", "Hi! How are you?", parentActorRef)
}

func TestForward(t *testing.T) {
	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor)

	forwarderActor := Actor{}
	defer forwarderActor.Close()
	forwarderActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
		forwarderActor.Forward(message, "childActor1", "childActor2")
	})
	ActorSystem().SpawnActor(&parentActor, "forwarderActor", &forwarderActor)

	childActor1 := Actor{}
	defer childActor1.Close()
	childActor1.React("message", func(message Message) {
		message.Self.LogInfo("Received %v from %v\n", message.Data, message.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor1", &childActor1)

	childActor2 := Actor{}
	defer childActor2.Close()
	childActor2.React("message", func(message Message) {
		message.Self.LogInfo("Received %v from %v\n", message.Data, message.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor2", &childActor2)

	parentActorRef, _ := ActorSystem().Actor("parentActor")
	forwarderActorRef, _ := ActorSystem().Actor("forwarderActor")

	forwarderActorRef.Send("message", "to be forwarded", parentActorRef)
	time.Sleep(500 * time.Millisecond)
}

func TestBecomeUnbecome(t *testing.T) {
	angry := func(message Message) {
		if message.Data == "happy" {
			message.Self.LogInfo("Unbecome\n")
			message.Self.Unbecome(message.messageType)
		} else {
			message.Self.LogInfo("Angrily receiving %v\n", message.Data)
		}
	}

	happy := func(message Message) {
		if message.Data == "angry" {
			message.Self.LogInfo("I shall become angry\n")
			message.Self.Become(message.messageType, angry)
		} else {
			message.Self.LogInfo("Happily receiving %v\n", message.Data)
		}
	}

	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor)

	childActor := Actor{}
	defer childActor.Close()
	childActor.React("message", happy)
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor)

	parentActorRef, _ := ActorSystem().Actor("parentActor")
	childActorRef, _ := ActorSystem().Actor("childActor")

	childActorRef.Send("message", "hello!", parentActorRef)
	childActorRef.Send("message", "angry", parentActorRef)
	childActorRef.Send("message", "hello!", parentActorRef)
	childActorRef.Send("message", "happy", parentActorRef)
	childActorRef.Send("message", "hello!", parentActorRef)
}
