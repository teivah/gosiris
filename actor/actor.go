package actor

type Actor struct {
	name    string
	conf    map[string]func(Message)
	mailbox chan Message
}

type actorInterface interface {
	React(string, func(Message))
	Configuration() map[string]func(Message)
	Mailbox() chan Message
	setMailbox(chan Message)
}

func (actor *Actor) Stringer() string {
	return actor.name
}

func (actor *Actor) React(messageType string, f func(Message)) {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Message))
	}

	actor.conf[messageType] = f
}

func (actor *Actor) Configuration() map[string]func(Message) {
	return actor.conf
}

func (actor *Actor) Mailbox() chan Message {
	return actor.mailbox
}

func (actor *Actor) setMailbox(mailbox chan Message) {
	actor.mailbox = mailbox
}
