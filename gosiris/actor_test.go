package gosiris

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

func TestBasic(t *testing.T) {
	t.Log("Starting Basic test")

	//Configure the actor system options
	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}

	//Init a local actor system
	InitActorSystem(opts)
	//Defer the actor system closure
	defer CloseActorSystem()

	//Create a simple parent actor
	parentActor := Actor{}
	//Defer the actor closure
	defer parentActor.Close()

	//Register the parent actor
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	//Create a simple child actor
	childActor := Actor{}
	//Defer the actor system closure
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

	//Tell a message from the parent to the child actor
	childActorRef.Tell("message", "Hi! How are you?", parentActorRef)

	time.Sleep(250 * time.Millisecond)
}

func TestStatefulness(t *testing.T) {
	t.Log("Starting statefulness test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

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
			message.Sender.Tell("error", "I already know you!", message.Self)
			childActor.Parent().Tell("help", "Daddy help me!", message.Self)
			childActor.Close()
		} else {
			childActor.hello[message.Data.(string)] = true
			message.Sender.Tell("helloback", "hello "+name+"!", message.Self)
		}
	})
	ActorSystem().SpawnActor(&parentActor, "child", &childActor, nil)

	childActorRef, _ := ActorSystem().ActorOf("child")
	parentActorRef, _ := ActorSystem().ActorOf("parent")

	childActorRef.Tell("hello", "teivah", parentActorRef)
	childActorRef.Tell("hello", "teivah", parentActorRef)

	time.Sleep(250 * time.Millisecond)
}

func TestForward(t *testing.T) {
	t.Log("Starting forward test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	forwarderActor := Actor{}
	defer forwarderActor.Close()
	forwarderActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
		message.Self.Forward(message, "childActor1", "childActor2")
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

	forwarderActorRef.Tell("message", "to be forwarded", parentActorRef)
	time.Sleep(1500 * time.Millisecond)
}

func TestBecomeUnbecome(t *testing.T) {
	t.Log("Starting become/unbecome test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	angry := func(message Message) {
		if message.Data == "happy" {
			message.Self.LogInfo("Unbecome\n")
			message.Self.Unbecome(message.MessageType)
		} else {
			message.Self.LogInfo("Angrily receiving %v\n", message.Data)
		}
	}

	happy := func(message Message) {
		if message.Data == "angry" {
			message.Self.LogInfo("I shall become angry\n")
			message.Self.Become(message.MessageType, angry)
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

	childActorRef.Tell("message", "hello!", parentActorRef)
	childActorRef.Tell("message", "angry", parentActorRef)
	childActorRef.Tell("message", "hello!", parentActorRef)
	childActorRef.Tell("message", "happy", parentActorRef)
	childActorRef.Tell("message", "hello!", parentActorRef)
	time.Sleep(500 * time.Millisecond)
}

func TestAmqp(t *testing.T) {
	t.Log("Starting AMQP test")
	time.Sleep(500 * time.Millisecond)

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actor1 := new(Actor).React("reply", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("actorX", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	actor2 := new(Actor).React("message", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
		message.Sender.Tell("reply", "hello back", message.Self)
	})
	defer actor2.Close()
	ActorSystem().RegisterActor("actorY", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor2"))

	actorRef1, _ := ActorSystem().ActorOf("actorX")
	actorRef2, _ := ActorSystem().ActorOf("actorY")

	actorRef2.Tell("message", "hello", actorRef1)
	time.Sleep(500 * time.Millisecond)
}

func TestKafka(t *testing.T) {
	t.Log("Starting Kafka test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actor1 := new(Actor).React("reply", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("actorX", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor1"))

	actor2 := new(Actor).React("message", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
		message.Sender.Tell("reply", "hello back", message.Self)
	})
	defer actor2.Close()
	ActorSystem().RegisterActor("actorY", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor2"))

	actorRef1, _ := ActorSystem().ActorOf("actorX")
	actorRef2, _ := ActorSystem().ActorOf("actorY")

	actorRef2.Tell("message", "hello", actorRef1)
	time.Sleep(2500 * time.Millisecond)
}

func TestAmqpKafka(t *testing.T) {
	t.Log("Starting AMQP/Kakfa test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
		ZipkinOptions: ZipkinOptions{
			Url:      "http://zipkin:9411/api/v1/spans",
			Debug:    true,
			HostPort: "0.0.0.0",
			SameSpan: true,
		},
	}

	InitActorSystem(opts)
	defer CloseActorSystem()

	actor1 := new(Actor).React("reply", func(message Message) {
		message.Self.LogInfo("Received1 %v", message.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("amqpActor", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	actor2 := new(Actor).React("message", func(message Message) {
		message.Self.LogInfo("Received2 %v", message.Data)
		message.Sender.Tell("reply", "hello back", message.Self)
	})
	defer actor2.Close()
	ActorSystem().SpawnActor(actor1, "kafkaActor", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor2"))

	amqpRef, _ := ActorSystem().ActorOf("amqpActor")
	kafkaRef, _ := ActorSystem().ActorOf("kafkaActor")

	//How to know when a tell is part of a span or not?
	kafkaRef.Tell("message", "hello", amqpRef)
	time.Sleep(1500 * time.Second)
}

func TestRemoteClose(t *testing.T) {
	t.Log("Starting remote close test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorY := new(Actor).React("hello", func(message Message) {
		message.Self.LogInfo("Received %v", message.Data)
		message.Sender.Tell("reply", fmt.Sprintf("Hello %v", message.Data), message.Self)
	})
	defer actorY.Close()
	ActorSystem().RegisterActor("actorY", actorY, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor2"))

	a, _ := ActorSystem().ActorOf("actorY")
	a.AskForClose(a)
	time.Sleep(500 * time.Millisecond)
}

func TestAutocloseTrue(t *testing.T) {
	t.Log("Starting auto close true test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorY := new(Actor)
	defer actorY.Close()
	ActorSystem().RegisterActor("actorY", actorY, new(ActorOptions).SetAutoclose(false))
	a, _ := ActorSystem().ActorOf("actorY")

	a.AskForClose(a)
	time.Sleep(500 * time.Millisecond)
}

func TestAutocloseFalse(t *testing.T) {
	t.Log("Starting auto close false test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorY := new(Actor)
	actorY.React(GosirisMsgPoisonPill, func(message Message) {
		message.Self.LogInfo("Received a poison pill, closing.")
		actorY.Close()
	})

	ActorSystem().RegisterActor("actorY", actorY, new(ActorOptions).SetAutoclose(false))
	a, _ := ActorSystem().ActorOf("actorY")

	a.AskForClose(a)
	time.Sleep(500 * time.Millisecond)
}

func TestChildClosedNotificationLocal(t *testing.T) {
	t.Log("Starting child closed notification test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorParent := new(Actor)
	defer actorParent.Close()
	actorParent.React(GosirisMsgChildClosed, func(message Message) {
		message.Self.LogInfo("My child is closed")
	})

	actorChild := new(Actor)
	actorChild.React("do", func(message Message) {
		if message.Data == 0 {
			message.Self.LogInfo("I feel like being closed")
			actorChild.Close()
		} else {
			message.Self.LogInfo("Received %v", message.Data)
		}
	})

	ActorSystem().RegisterActor("ActorParent", actorParent, nil)
	ActorSystem().SpawnActor(actorParent, "ActorChild", actorChild, nil)

	actorParentRef, _ := ActorSystem().ActorOf("ActorParent")
	actorChildRef, _ := ActorSystem().ActorOf("ActorChild")

	actorChildRef.Tell("do", 1, actorParentRef)
	actorChildRef.Tell("do", 0, actorParentRef)

	time.Sleep(500 * time.Millisecond)
}

func TestChildClosedNotificationRemote(t *testing.T) {
	t.Log("Starting child closed notification remote test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorParent := new(Actor)
	defer actorParent.Close()
	actorParent.React(GosirisMsgChildClosed, func(message Message) {
		message.Self.LogInfo("My child is closed")
	})

	actorChild := new(Actor)
	actorChild.React("do", func(message Message) {
		if message.Data == "0" {
			message.Self.LogInfo("I feel like being closed")
			actorChild.Close()
		} else {
			message.Self.LogInfo("Received %v", message.Data)
		}
	})

	ActorSystem().RegisterActor("ActorParent", actorParent, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actorParent"))
	ActorSystem().SpawnActor(actorParent, "ActorChild", actorChild, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actorChild"))

	actorParentRef, _ := ActorSystem().ActorOf("ActorParent")
	actorChildRef, _ := ActorSystem().ActorOf("ActorChild")

	actorChildRef.Tell("do", 1, actorParentRef)
	actorChildRef.Tell("do", 0, actorParentRef)

	time.Sleep(1000 * time.Millisecond)
}

func TestRepeat(t *testing.T) {
	t.Log("Starting repeat test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	parentActor := Actor{}
	defer parentActor.Close()

	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	childActor := Actor{}
	defer childActor.Close()

	childActor.React("message", func(message Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})

	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	c, _ := childActorRef.Repeat("message", 5*time.Millisecond, "Hi! How are you?", parentActorRef)

	time.Sleep(21 * time.Millisecond)

	ActorSystem().Stop(c)
	time.Sleep(500 * time.Millisecond)
}

//TODO
func TestDefaultWatcher(t *testing.T) {
	t.Log("Starting default watcher test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	parentActor := Actor{}
	defer parentActor.Close()

	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	childActor := Actor{}
	defer childActor.Close()

	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, new(ActorOptions).SetDefaultWatcher(5*time.Millisecond))

	time.Sleep(21 * time.Millisecond)
}
