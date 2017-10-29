package gofast

import (
	"fmt"
)

var root = &Actor{}

const (
	remote_amqp = "remote_amqp"
)

func RootActor() *Actor {
	return root
}

type Actor struct {
	name     string
	conf     map[string]func(Message)
	mailbox  chan Message
	parent   actorInterface
	unbecome map[string]func(Message)
}

type RemoteActor struct {
	Actor
}

type actorInterface interface {
	React(string, func(Message)) *Actor
	reactions() map[string]func(Message)
	Unbecome() map[string]func(Message)
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
	if actor.Mailbox() != nil { //If local
		close(actor.mailbox)
	}
	ActorSystem().unregisterActor(actor.name)
}

func (actor *Actor) Stringer() string {
	return actor.name
}

func (actor *Actor) React(messageType string, f func(Message)) *Actor {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Message))
		actor.unbecome = make(map[string]func(Message))
	}

	actor.conf[messageType] = f

	return actor
}

func (actor *Actor) reactions() map[string]func(Message) {
	return actor.conf
}

func (actor *Actor) Unbecome() map[string]func(Message) {
	return actor.unbecome
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
	parent, _ := ActorSystem().ActorOf(actor.parent.Name())
	return parent
}

func (actor *Actor) Name() string {
	return actor.name
}

func (actor *Actor) Forward(message Message, destinations ...string) {
	for _, v := range destinations {
		actorRef, err := ActorSystem().ActorOf(v)
		if err != nil {
			fmt.Errorf("actor %v is not part of the actor system", v)
		}
		actorRef.Send(message.messageType, message.Data, message.Sender)
	}
}
