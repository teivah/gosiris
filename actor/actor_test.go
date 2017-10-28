package actor

import (
	"testing"
	"Test/src/actor"
	"fmt"
	"time"
)

func TestRequestReply(t *testing.T) {
	child := &actor.Actor{}
	child.React("hello", func(message actor.Message) {
		fmt.Printf("Receive %v\n", message.Data)
		message.Sender.Tell("response", "hello teivah!", message.Self)
	})
	child.Build("child")

	parent := &actor.Actor{}
	parent.React("response", func(message actor.Message) {
		fmt.Printf("Receive response %v\n", message.Data)
	})
	parent.Build("parent")

	err := child.Tell("hello", "teivah", parent)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(1000)
}
