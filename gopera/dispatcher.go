package gopera

import (
	"gopera/gopera/util"
	"encoding/json"
	"fmt"
)

func init() {

}

const (
	json_messageType = "messageType"
	json_data        = "data"
	json_sender      = "sender"
	json_self        = "self"
)

type Message struct {
	messageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
}

func (message Message) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m[json_messageType] = message.messageType
	m[json_data] = fmt.Sprint(message.Data)
	m[json_sender] = message.Sender.Name()
	m[json_self] = message.Self.Name()
	return json.Marshal(m)
}

func (message *Message) UnmarshalJSON(b []byte) error {
	var m map[string]string
	err := json.Unmarshal(b, &m)
	if err != nil {
		util.LogError("Unmarshalling error: %v", err)
		return err
	}

	message.messageType = m[json_messageType]

	message.Data = m[json_data]

	self := m[json_self]
	selfAssociation, err := ActorSystem().actor(self)
	if err != nil {
		return err
	}
	message.Self = selfAssociation.actorRef

	sender := m[json_sender]
	senderAssociation, err := ActorSystem().actor(sender)
	if err != nil {
		return err
	}
	message.Sender = senderAssociation.actorRef

	return nil
}

var poisonPill = "poisonpill"

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface, options OptionsInterface) error {
	defer func() {
		if r := recover(); r != nil {
			util.LogInfo("Recovered in %v", r)
		}
	}()

	m := Message{messageType, data, sender, receiver}

	if !options.Remote() {
		channel <- m
		util.LogInfo("Message dispatched to local channel")
	} else {
		d, err := RemoteConnection(receiver.Name())
		if err != nil {
			return err
		}

		json, err := json.Marshal(m)
		if err != nil {
			util.LogError("JSON marshalling error: %v", err)
			return err
		}

		d.Send(options.Destination(), json)
		util.LogInfo("Message dispatched to remote channel %v", options.Destination())
	}

	return nil
}

func receive(actor actorInterface, options OptionsInterface) {
	if !options.Remote() {
		defer func() {
			if r := recover(); r != nil {
				util.LogInfo("Recovered in %v", r)
			}
		}()

		c := actor.Mailbox()
		for {
			select {
			case p := <-c:
				ActorSystem().Invoke(p)
			}
		}
	} else {
		d, err := RemoteConnection(actor.Name())
		if err != nil {
			return
		}

		d.Receive(options.Destination())
	}
}
