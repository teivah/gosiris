package gofast

var poisonPill Message

func PoisonPill() Message {
	return poisonPill
}

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
		select {
		case p := <-c:
			if p == poisonPill {
				close(c)
				return
			}
			actor.configuration()[p.messageType](p)
		}
	}
}
