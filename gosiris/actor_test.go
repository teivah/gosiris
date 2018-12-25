package gosiris

import (
	"fmt"
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
	t.Log("Starting Basic test")

	//Init a local actor system
	InitActorSystem(SystemOptions{
		ActorSystemName: "ActorSystem",
	})
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

	//Register the reactions to event types (here a reaction to context)
	childActor.React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v\n", context.Data)
	})

	//Register the child actor
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	//Retrieve the parent and child actor reference
	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	//Tell a context from the parent to the child actor
	childActorRef.Tell(EmptyContext, "context", "Hi! How are you?", parentActorRef)

	time.Sleep(2500 * time.Millisecond)
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

	f := func(context Context) {
		context.Self.LogInfo(context, "Receive response %v\n", context.Data)
	}

	parentActor.React("helloback", f).React("error", f).React("help", f)
	ActorSystem().RegisterActor("parent", &parentActor, nil)

	childActor.React("hello", func(context Context) {
		context.Self.LogInfo(context, "Receive request %v\n", context.Data)

		name := context.Data.(string)

		if _, ok := childActor.hello[name]; ok {
			context.Sender.Tell(context, "error", "I already know you!", context.Self)
			childActor.Parent().Tell(context, "help", "Daddy help me!", context.Self)
			childActor.Close()
		} else {
			childActor.hello[context.Data.(string)] = true
			context.Sender.Tell(context, "helloback", "hello "+name+"!", context.Self)
		}
	})
	ActorSystem().SpawnActor(&parentActor, "child", &childActor, nil)

	childActorRef, _ := ActorSystem().ActorOf("child")
	parentActorRef, _ := ActorSystem().ActorOf("parent")

	childActorRef.Tell(EmptyContext, "hello", "teivah", parentActorRef)
	childActorRef.Tell(EmptyContext, "hello", "teivah", parentActorRef)

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
	forwarderActor.React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v\n", context.Data)
		context.Self.Forward(context, "childActor1", "childActor2")
	})
	ActorSystem().SpawnActor(&parentActor, "forwarderActor", &forwarderActor, nil)

	childActor1 := Actor{}
	defer childActor1.Close()
	childActor1.React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v from %v\n", context.Data, context.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor1", &childActor1, nil)

	childActor2 := Actor{}
	defer childActor2.Close()
	childActor2.React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v from %v\n", context.Data, context.Sender)
	})
	ActorSystem().SpawnActor(&forwarderActor, "childActor2", &childActor2, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	forwarderActorRef, _ := ActorSystem().ActorOf("forwarderActor")

	forwarderActorRef.Tell(EmptyContext, "context", "to be forwarded", parentActorRef)
	time.Sleep(1500 * time.Millisecond)
}

func TestBecomeUnbecome(t *testing.T) {
	t.Log("Starting become/unbecome test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	angry := func(context Context) {
		if context.Data == "happy" {
			context.Self.LogInfo(context, "Unbecome\n")
			context.Self.Unbecome(context.MessageType)
		} else {
			context.Self.LogInfo(context, "Angrily receiving %v\n", context.Data)
		}
	}

	happy := func(context Context) {
		if context.Data == "angry" {
			context.Self.LogInfo(context, "I shall become angry\n")
			context.Self.Become(context.MessageType, angry)
		} else {
			context.Self.LogInfo(context, "Happily receiving %v\n", context.Data)
		}
	}

	parentActor := Actor{}
	defer parentActor.Close()
	ActorSystem().RegisterActor("parentActor", &parentActor, nil)

	childActor := Actor{}
	defer childActor.Close()
	childActor.React("context", happy)
	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	childActorRef.Tell(EmptyContext, "context", "hello!", parentActorRef)
	childActorRef.Tell(EmptyContext, "context", "angry", parentActorRef)
	childActorRef.Tell(EmptyContext, "context", "hello!", parentActorRef)
	childActorRef.Tell(EmptyContext, "context", "happy", parentActorRef)
	childActorRef.Tell(EmptyContext, "context", "hello!", parentActorRef)
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

	actor1 := new(Actor).React("reply", func(context Context) {
		context.Self.LogInfo(context, "Received %v", context.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("actorX", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	actor2 := new(Actor).React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v", context.Data)
		context.Sender.Tell(context, "reply", "hello back", context.Self)
	})
	defer actor2.Close()
	ActorSystem().RegisterActor("actorY", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor2"))

	actorRef1, _ := ActorSystem().ActorOf("actorX")
	actorRef2, _ := ActorSystem().ActorOf("actorY")

	actorRef2.Tell(EmptyContext, "context", "hello", actorRef1)
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

	actor1 := new(Actor).React("reply", func(context Context) {
		context.Self.LogInfo(context, "Received %v", context.Data)
	})
	defer actor1.Close()
	ActorSystem().RegisterActor("actorX", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor1"))

	actor2 := new(Actor).React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v", context.Data)
		context.Sender.Tell(EmptyContext, "reply", "hello back", context.Self)
	})
	defer actor2.Close()
	ActorSystem().RegisterActor("actorY", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor2"))

	actorRef1, _ := ActorSystem().ActorOf("actorX")
	actorRef2, _ := ActorSystem().ActorOf("actorY")

	actorRef2.Tell(EmptyContext, "context", "hello", actorRef1)
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

	actor1 := new(Actor).React("reply", func(context Context) {
		context.Self.LogInfo(context, "Received1 %v", context.Data)

	})
	defer actor1.Close()
	ActorSystem().RegisterActor("amqpActor", actor1, new(ActorOptions).SetRemote(true).SetRemoteType(Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	actor2 := new(Actor).React("context", func(context Context) {
		context.Self.LogInfo(context, "Received2 %v", context.Data)
		context.Sender.Tell(context, "reply", "hello back", context.Self)
	})
	defer actor2.Close()
	ActorSystem().SpawnActor(actor1, "kafkaActor", actor2, new(ActorOptions).SetRemote(true).SetRemoteType(Kafka).SetUrl("kafka:9092").SetDestination("actor2"))

	amqpRef, _ := ActorSystem().ActorOf("amqpActor")
	kafkaRef, _ := ActorSystem().ActorOf("kafkaActor")

	kafkaRef.Tell(EmptyContext, "context", "hello", amqpRef)
	time.Sleep(1500 * time.Millisecond)
}

func TestRemoteClose(t *testing.T) {
	t.Log("Starting remote close test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	actorY := new(Actor).React("hello", func(context Context) {
		context.Self.LogInfo(context, "Received %v", context.Data)
		context.Sender.Tell(context, "reply", fmt.Sprintf("Hello %v", context.Data), context.Self)
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
	actorY.React(GosirisMsgPoisonPill, func(context Context) {
		context.Self.LogInfo(context, "Received a poison pill, closing.")
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
	actorParent.React(GosirisMsgChildClosed, func(context Context) {
		context.Self.LogInfo(context, "My child is closed")
	})

	actorChild := new(Actor)
	actorChild.React("do", func(context Context) {
		if context.Data == 0 {
			context.Self.LogInfo(context, "I feel like being closed")
			actorChild.Close()
		} else {
			context.Self.LogInfo(context, "Received %v", context.Data)
		}
	})

	ActorSystem().RegisterActor("ActorParent", actorParent, nil)
	ActorSystem().SpawnActor(actorParent, "ActorChild", actorChild, nil)

	actorParentRef, _ := ActorSystem().ActorOf("ActorParent")
	actorChildRef, _ := ActorSystem().ActorOf("ActorChild")

	actorChildRef.Tell(EmptyContext, "do", 1, actorParentRef)
	actorChildRef.Tell(EmptyContext, "do", 0, actorParentRef)

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
	actorParent.React(GosirisMsgChildClosed, func(context Context) {
		context.Self.LogInfo(context, "My child is closed")
	})

	actorChild := new(Actor)
	actorChild.React("do", func(context Context) {
		if context.Data == "0" {
			context.Self.LogInfo(context, "I feel like being closed")
			actorChild.Close()
		} else {
			context.Self.LogInfo(context, "Received %v", context.Data)
		}
	})

	ActorSystem().RegisterActor("ActorParent", actorParent, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actorParent"))
	ActorSystem().SpawnActor(actorParent, "ActorChild", actorChild, new(ActorOptions).SetRemote(true).SetRemoteType("amqp").SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actorChild"))

	actorParentRef, _ := ActorSystem().ActorOf("ActorParent")
	actorChildRef, _ := ActorSystem().ActorOf("ActorChild")

	actorChildRef.Tell(EmptyContext, "do", 1, actorParentRef)
	actorChildRef.Tell(EmptyContext, "do", 0, actorParentRef)

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

	childActor.React("context", func(context Context) {
		context.Self.LogInfo(context, "Received %v\n", context.Data)
	})

	ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	parentActorRef, _ := ActorSystem().ActorOf("parentActor")
	childActorRef, _ := ActorSystem().ActorOf("childActor")

	c, _ := childActorRef.Repeat("context", 5*time.Millisecond, "Hi! How are you?", parentActorRef)

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

func TestAsk(t *testing.T) {
	t.Log("Starting ask test")

	opts := SystemOptions{
		ActorSystemName: "ActorSystem",
	}
	InitActorSystem(opts)
	defer CloseActorSystem()

	replierActor := Actor{}
	defer replierActor.Close()

	replierActor.React("time", func(context Context) {
		context.Sender.Tell(EmptyContext, "time", time.Now(), context.Self)
	})

	ActorSystem().RegisterActor("replierActor", &replierActor, nil)

	replierActorRef, _ := ActorSystem().ActorOf("replierActor")

	successfulReply, err := replierActorRef.Ask(EmptyContext, "time", nil, 1 * time.Second)
	if err != nil {
		t.Error("reply not received")
		t.Fail()
	}
	t.Logf("successfulReply was: %v", successfulReply)

	reply, timeoutError := replierActorRef.Ask(EmptyContext, "unknownQuestion", nil, 1 * time.Second)
	if reply != nil {
		t.Error("this should have timed out")
		t.Fail()
	}
	t.Logf("unknownQuestion timedout: %v", timeoutError.Error())
}
