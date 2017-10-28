package actor

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
	for {
		p := <-c
		actor.configuration()[p.messageType](p)
	}
}
