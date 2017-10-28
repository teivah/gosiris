package actor

import (
	"testing"
	"fmt"
	"time"
)

type ParentActor struct {
	Actor
}

type ChildActor struct {
	Actor
	hello map[string]bool
}

var parentActor ParentActor
var childActor ChildActor

func init() {
	childActor.hello = make(map[string]bool)
}

func TestStatefulness(t *testing.T) {
	childActor.React("hello", func(message Message) {
		fmt.Printf("Child: receive %v\n", message.Data)

		name := message.Data.(string)

		if _, ok := childActor.hello[name]; ok {
			message.Sender.Tell("error", "I already know you!", message.Self)
		} else {
			childActor.hello[message.Data.(string)] = true
			message.Sender.Tell("helloback", "hello "+name+"!", message.Self)
		}
	})
	ActorSystem().RegisterActor("child", &childActor)

	parentActor.React("helloback", func(message Message) {
		fmt.Printf("Parent: receive response %v\n", message.Data)
	})
	parentActor.React("error", func(message Message) {
		fmt.Printf("Parent: receive response %v\n", message.Data)
	})
	ActorSystem().RegisterActor("parent", &parentActor)

	childActorRef, _ := ActorSystem().Actor("child")
	parentActorRef, _ := ActorSystem().Actor("parent")

	childActorRef.Tell("hello", "teivah", parentActorRef)
	childActorRef.Tell("hello", "teivah", parentActorRef)

	time.Sleep(1000)
}
