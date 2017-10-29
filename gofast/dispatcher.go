package gofast

import (
	"flogo-lib/logger"
	"Gofast/gofast/util"
	"encoding/json"
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

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface, options OptionsInterface) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Info("Recovered in %v", r)
		}
	}()

	m := Message{messageType, data, sender, receiver}

	if options == nil {
		channel <- m
		util.LogInfo("Message dispatched to local channel")
	} else {
		d := ActorSystem().DistributedConfiguration(options.ConnectionAlias())
		if d == nil {
			util.LogError("remote configuration %v cannot be found", options.ConnectionAlias())
		}

		json, err := json.Marshal(m)
		if err != nil {
			util.LogError("JSON marshalling error: %v", err)
			return err
		}

		d.Send(options.Endpoint(), json)
		util.LogInfo("Message dispatched to remote channel %v", options.Endpoint())
	}

	return nil
}

func receive(actor actorInterface, options OptionsInterface) {
	if options == nil {
		defer func() {
			if r := recover(); r != nil {
				logger.Info("Recovered in %v", r)
			}
		}()

		c := actor.Mailbox()
		for {
			select {
			case p := <-c:
				if p.messageType == poisonPill {
					util.LogInfo("Actor %v has received a poison pill", actor.Name())
					actor.Close()
					return
				}

				f, exists := actor.reactions()[p.messageType]
				if exists {
					f(p)
				}
			}
		}
	} else {
		d := ActorSystem().DistributedConfiguration(options.ConnectionAlias())
		if d == nil {
			util.LogError("remote configuration %v cannot be found", options.ConnectionAlias())
		}

		d.Receive(options.Endpoint())
	}
}
