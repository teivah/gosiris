package actor

import "fmt"

type Actor struct {
	name    string
	conf    map[string]func(Message)
	built   bool
	mailbox chan Message
}

type actorInterface interface {
	Build(string)
	React(string, func(Message))
	Tell(string, interface{}, actorInterface) error
	Configuration() map[string]func(Message)
	Name() string
	Channel() chan Message
}

func (actor *Actor) Build(name string) {
	actor.name = name
	actor.mailbox = make(chan Message)

	go receive(actor)

	actor.built = true
}

func (actor *Actor) React(messageType string, f func(Message)) {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Message))
	}

	actor.conf[messageType] = f
}

func (actor *Actor) Tell(messageType string, data interface{}, sender actorInterface) error {
	if !actor.built {
		return fmt.Errorf("actor is not configured")
	}

	if actor.conf == nil || actor.conf[messageType] == nil {
		return fmt.Errorf("actor %v is not configured to receive %v", actor, messageType)
	}

	dispatch(messageType, data, actor, sender)

	return nil
}

func (actor *Actor) Configuration() map[string]func(Message) {
	return actor.conf
}

func (actor *Actor) Name() string {
	return actor.name
}

func (actor *Actor) Stringer() string {
	return actor.name
}

func (actor *Actor) Channel() chan Message {
	return actor.mailbox
}
