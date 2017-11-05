package gopera

import (
	"encoding/json"
	"fmt"
)

func init() {

}

const (
	GoperaMsgPoisonPill  = "GoperaPoisonPill"
	GoperaMsgChildClosed = "GoperaChildClosed"
	jsonMessageType      = "messageType"
	jsonData             = "data"
	jsonSender           = "sender"
	jsonSelf             = "self"
)

type Message struct {
	MessageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
}

func (message Message) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m[jsonMessageType] = message.MessageType
	m[jsonData] = fmt.Sprint(message.Data)
	m[jsonSender] = message.Sender.Name()
	m[jsonSelf] = message.Self.Name()
	return json.Marshal(m)
}

func (message *Message) UnmarshalJSON(b []byte) error {
	var m map[string]string
	err := json.Unmarshal(b, &m)
	if err != nil {
		ErrorLogger.Printf("Unmarshalling error: %v", err)
		return err
	}

	message.MessageType = m[jsonMessageType]

	message.Data = m[jsonData]

	self := m[jsonSelf]
	selfAssociation, err := ActorSystem().actor(self)
	if err != nil {
		return err
	}
	message.Self = selfAssociation.actorRef

	sender := m[jsonSender]
	senderAssociation, err := ActorSystem().actor(sender)
	if err != nil {
		return err
	}
	message.Sender = senderAssociation.actorRef

	return nil
}

func dispatch(channel chan Message, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface, options OptionsInterface) error {
	defer func() {
		if r := recover(); r != nil {
			InfoLogger.Printf("Dispatch recovered in %v", r)
		}
	}()

	InfoLogger.Printf("Dispatching message %v", messageType)

	m := Message{messageType, data, sender, receiver}

	if !options.Remote() {
		channel <- m
		InfoLogger.Printf("Message dispatched to local channel")
	} else {
		d, err := RemoteConnection(receiver.Name())
		if err != nil {
			return err
		}

		json, err := json.Marshal(m)
		if err != nil {
			ErrorLogger.Printf("JSON marshalling error: %v", err)
			return err
		}

		d.Send(options.Destination(), json)
		InfoLogger.Printf("Message dispatched to remote channel %v", options.Destination())
	}

	return nil
}

func receive(actor actorInterface, options OptionsInterface) {
	if !options.Remote() {
		defer func() {
			if r := recover(); r != nil {
				InfoLogger.Printf("Receive recovered in %v", r)
			}
		}()

		c := actor.getDataChan()
		for {
			select {
			//TODO To be improved. If the channel is closed it will trigger the recover function
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
