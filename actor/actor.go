package actor

import (
	"fmt"
)

type Actor struct {
	name    string
	conf    map[string]func(Message)
	mailbox chan Message
}

type actorInterface interface {
	Printf(string, ...interface{}) (int, error)
	React(string, func(Message))
	configuration() map[string]func(Message)
	Mailbox() chan Message
	setMailbox(chan Message)
	setName(string)
}

func (actor *Actor) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("["+actor.name+"] "+format, a)
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

func (actor *Actor) configuration() map[string]func(Message) {
	return actor.conf
}

func (actor *Actor) Mailbox() chan Message {
	return actor.mailbox
}

func (actor *Actor) setMailbox(mailbox chan Message) {
	actor.mailbox = mailbox
}

func (actor *Actor) setName(name string) {
	actor.name = name
}
