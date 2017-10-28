package gofast

import (
	"fmt"
)

type Message struct {
	messageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
}

var poisonPill = Message{"poisonpill", nil, nil, nil}

func PoisonPill() Message {
	return poisonPill
}

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered2 in f", r)
		}
	}()

	channel <- Message{messageType, data, sender, receiver}
}

func receive(actor actorInterface) {
	c := actor.Mailbox()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered1 in f", r)
		}
	}()

	for {
		select {
		case p := <-c:
			if p == PoisonPill() {
				actor.Close()
				return
			}

			f, exists := actor.reactions()[p.messageType]
			if exists {
				f(p)
			}
		}
	}
}
