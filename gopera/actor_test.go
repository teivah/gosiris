package gopera

import (
	"testing"
	"time"
	"fmt"
)

type ParentActor struct {
	Actor
}

type ChildActor struct {
	Actor
	hello map[string]bool
}

func init() {
	InitDistributedActorSystem("http://etcd:2379")
}

func TestBasic(t *testing.T) {
	//Create a simple parent actor
	parentActor := Actor{}
	//Close the actor
	defer parentActor.Close()

	//Register the parent actor
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	//Create a simple child actor
	childActor := Actor{}
	//Close the actor
	defer childActor.Close()

	//Register the reactions to event types (here a reaction to message)
	childActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})

	//Register the child actor
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	//Retrieve the parent and child actor reference
	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

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
	ActorSystem().RegisterActor("parent", &parentActor, nil)

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
	ActorSystem().SpawnActor(&parentActor, "child", &childActor, nil)

	childActorRef, _ := ActorSystem().ActorOf("child")
	parentActorRef, _ := ActorSystem().ActorOf("parent")

	childActorRef.Send("hello", "teivah", parentActorRef)
	childActorRef.Send("hello", "teivah", parentActorRef)

	time.Sleep(500 * time.Millisecond)
	childActorRef.Send("hello", "teivah2", parentActorRef)
}

func TestClose(t *testing.T) {
	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	childActor := Actor{}
	childActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	childActorRef.AskForClose(parentActorRef)

	time.Sleep(500 * time.Millisecond)
	childActorRef.Send("message", "Hi! How are you?", parentActorRef)
}

func TestForward(t *testing.T) {
	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	forwarderActor := Actor{}
	defer forwarderActor.Close()
	forwarderActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
		forwarderActor.Forward(message, "childActor1", "childActor2")
	})
	ActorSystem().SpawnActor(&parentActor, "forwarderActor", &forwarderActor, nil)

	childActor1 := Actor{}
	defer childActor1.Close()
	childActor1.React("message", func(message Message) {
		message.Self.LogInfo("Received %v from %v\n", message.Data, message.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor1", &childActor1, nil)

	childActor2 := Actor{}
	defer childActor2.Close()
	childActor2.React("message", func(message Message) {
		message.Self.LogInfo("Received %v from %v\n", message.Data, message.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor2", &childActor2, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	forwarderActorRef, _ := ActorSystem().ActorOf("forwarderActor")

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
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	childActor := Actor{}
	defer childActor.Close()
	childActor.React("message", happy)
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	childActorRef.Send("message", "hello!", parentActorRef)
	childActorRef.Send("message", "angry", parentActorRef)
	childActorRef.Send("message", "hello!", parentActorRef)
	childActorRef.Send("message", "happy", parentActorRef)
	childActorRef.Send("message", "hello!", parentActorRef)
}

func TestNewRemoteActor(t *testing.T) {
	actor1 := new(Actor).React("reply", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("actorX", actor1, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	actor2 := new(Actor).React("message", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
		message.Sender.Send("reply", "hello back", message.Self)
	})
	defer actor2.Close()
	ActorSystem().RegisterActor("actorY", actor2, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor2"))

	actorRef1, _ := ActorSystem().ActorOf("actorX")
	actorRef2, _ := ActorSystem().ActorOf("actorY")

	actorRef2.Send("message", "hello", actorRef1)
	time.Sleep(500 * time.Millisecond)
}

func TestUnregister(t *testing.T) {
	actorY := new(Actor).React("hello", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
		message.Sender.Send("reply", fmt.Sprintf("Hello %v", message.Data), message.Self)
	})
	defer actorY.Close()
	ActorSystem().RegisterActor("actorY", actorY, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor2"))

	a, _ := ActorSystem().ActorOf("actorY")
	a.AskForClose(a)
}
