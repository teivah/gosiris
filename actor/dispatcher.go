package actor

type Message struct {
	messageType string
	Data        interface{}
	Sender      actorInterface
	Self        actorInterface
}

func dispatch(messageType string, data interface{}, actor actorInterface, sender actorInterface) {
	actor.Channel() <- Message{messageType, data, sender, actor}
}

func receive(actor actorInterface) {
	c := actor.Channel()
	p := <-c

	actor.Configuration()[p.messageType](p)
}