package gofast

type Message struct {
	messageType string
	Data        interface{}
	Sender      actorRefInterface
	Self        actorRefInterface
}

func dispatch(channel chan Message, messageType string, data interface{}, receiver actorRefInterface, sender actorRefInterface) {
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
			actor.configuration()[p.messageType](p)
		}
	}
}
