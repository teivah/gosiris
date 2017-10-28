package gofast

import (
	"fmt"
)

var root = &Actor{}

func RootActor() *Actor {
	return root
}

type Actor struct {
	name    string
	conf    map[string]func(Message)
	mailbox chan Message
	parent  actorInterface
}

type actorInterface interface {
	Printf(string, ...interface{}) (int, error)
	React(string, func(Message)) *Actor
	configuration() map[string]func(Message)
	Mailbox() chan Message
	setMailbox(chan Message)
	setName(string)
	setParent(actorInterface)
	Parent() ActorRefInterface
	Name() string
	Close()
	Forward(Message, ...string)
}

func (actor *Actor) Close() {
	close(actor.mailbox)
}

func (actor *Actor) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("["+actor.name+"] "+format, a)
}

func (actor *Actor) Stringer() string {
	return actor.name
}

func (actor *Actor) React(messageType string, f func(Message)) *Actor {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Message))
	}

	actor.conf[messageType] = f

	return actor
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

func (actor *Actor) setParent(parent actorInterface) {
	actor.parent = parent
}

func (actor *Actor) Parent() ActorRefInterface {
	parent, _ := ActorSystem().Actor(actor.parent.Name())
	return parent
}

func (actor *Actor) Name() string {
	return actor.name
}

func (actor *Actor) Forward(message Message, destinations ...string) {
	for _, v := range destinations {
		actorRef, err := ActorSystem().Actor(v)
		if err != nil {
			fmt.Errorf("actor %v is not part of the actor system", v)
		}
		actorRef.Send(message.messageType, message.Data, message.Sender)
	}
}
