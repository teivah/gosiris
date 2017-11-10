package gosiris

import (
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
)

func init() {

}

var (
	EmptyContext Context = Context{}
)

const (
	GosirisMsgPoisonPill       = "gosirisPoisonPill"
	GosirisMsgChildClosed      = "gosirisChildClosed"
	GosirisMsgHeartbeatRequest = "gosirisHeartbeatRequest"
	GosirisMsgHeartbeatReply   = "gosirisHeartbeatReply"

	jsonMessageType = "messageType"
	jsonData        = "data"
	jsonSender      = "sender"
	jsonSelf        = "self"
	jsonTracing     = "tracing"
)

type Context struct {
	MessageType string
	Data        interface{}
	Sender      ActorRefInterface
	Self        ActorRefInterface
	carrier     opentracing.TextMapCarrier
	span        opentracing.Span
}

func (context Context) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m[jsonMessageType] = context.MessageType
	m[jsonData] = fmt.Sprint(context.Data)
	m[jsonSender] = context.Sender.Name()
	m[jsonSelf] = context.Self.Name()
	m[jsonTracing] = context.carrier
	return json.Marshal(m)
}

func (context *Context) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		ErrorLogger.Printf("Unmarshalling error: %v", err)
		return err
	}

	context.MessageType = m[jsonMessageType].(string)

	context.Data = m[jsonData]

	self := m[jsonSelf].(string)
	selfAssociation, err := ActorSystem().actor(self)
	if err != nil {
		return err
	}
	context.Self = selfAssociation.actorRef

	sender := m[jsonSender].(string)
	senderAssociation, err := ActorSystem().actor(sender)
	if err != nil {
		return err
	}
	context.Sender = senderAssociation.actorRef

	if value, exists := m[jsonTracing]; exists {
		if value != nil {
			t := value.(map[string]interface{})
			context.carrier = make(map[string]string)

			for k, v := range t {
				context.carrier[k] = v.(string)
			}
		}
	}

	return nil
}

func dispatch(channel chan Context, messageType string, data interface{}, receiver ActorRefInterface, sender ActorRefInterface, options OptionsInterface, span opentracing.Span) error {
	defer func() {
		if r := recover(); r != nil {
			ErrorLogger.Printf("Dispatch recovered in %v", r)
		}
	}()

	InfoLogger.Printf("Dispatching message %v from %v to %v", messageType, sender.Name(), receiver.Name())

	carrier, _ := inject(span)
	m := Context{messageType, data, sender, receiver, carrier, nil}

	if !options.Remote() {
		channel <- m
		InfoLogger.Printf("Context dispatched to local channel")
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
		InfoLogger.Printf("Context dispatched to remote channel %v", options.Destination())
	}

	return nil
}

func receive(actor actorInterface, options OptionsInterface) {
	if !options.Remote() {
		defer func() {
			if r := recover(); r != nil {
				ErrorLogger.Printf("Receive recovered in %v", r)
			}
		}()

		dataChan := actor.getDataChan()
		closeChan := actor.getCloseChan()
		for {
			select {
			case p := <-dataChan:
				ActorSystem().Invoke(p)
			case <-closeChan:
				InfoLogger.Printf("Closing %v receiver", actor.Name())
				close(dataChan)
				close(closeChan)
				return
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
