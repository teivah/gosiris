package gofast

import (
	"flogo-lib/logger"
)

func init() {

}

type Message struct {
	messageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
}

var poisonPill = "poisonpill"

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface) {
	defer func() {
		if r := recover(); r != nil {
			logger.Info("Recovered in %v", r)
		}
	}()

	channel <- Message{messageType, data, sender, receiver}
}

func receive(actor actorInterface) {
	c := actor.Mailbox()

	defer func() {
		if r := recover(); r != nil {
			logger.Info("Recovered in %v", r)
		}
	}()

	for {
		select {
		case p := <-c:
			if p.messageType == poisonPill {
				logger.Info("Actor %v has received a poison pill", actor.Name())
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
