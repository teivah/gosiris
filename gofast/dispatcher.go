package gofast

type Message struct {
	messageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
}

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	channel <- Message{messageType, data, sender, receiver}
}

func receive(actor actorInterface) {
	c := actor.Mailbox()

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	for {
		select {
		case p := <-c:
			//fmt.Printf("receive %v\n", p)
			actor.configuration()[p.messageType](p)
		}
	}
}
