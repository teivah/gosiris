package gofast

type ActorRef struct {
	name string
}

type ActorRefInterface interface {
	Send(string, interface{}, ActorRefInterface) error
	AskForClose(ActorRefInterface)
}

func (ref ActorRef) Send(messageType string, data interface{}, sender ActorRefInterface) error {
	actor, _ := ActorSystem().actor(ref)

	dispatch(actor.Mailbox(), messageType, data, &ref, sender)

	return nil
}

func (ref ActorRef) AskForClose(sender ActorRefInterface) {
	actor, err := ActorSystem().actor(ref)
	
	if err != nil {
		return
	}

	dispatch(actor.Mailbox(), poisonPill.messageType, poisonPill.Data, &ref, sender)
}
