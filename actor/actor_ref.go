package actor

type ActorRef struct {
	name string
}

type actorRefInterface interface {
	Tell(string, interface{}, actorRefInterface) error
}

func (ref ActorRef) Tell(messageType string, data interface{}, sender actorRefInterface) error {
	actor, _ := ActorSystem().actor(ref)

	dispatch(actor.Mailbox(), messageType, data, &ref, sender)

	return nil
}
