package gofast

type ActorRef struct {
	name string
}

type ActorRefInterface interface {
	Tell(string, interface{}, ActorRefInterface) error
}

func (ref ActorRef) Tell(messageType string, data interface{}, sender ActorRefInterface) error {
	actor, _ := ActorSystem().actor(ref)

	dispatch(actor.Mailbox(), messageType, data, &ref, sender)

	return nil
}
