package actor

import (
	"fmt"
)

var actorSystemInstance actorSystem = actorSystem{}

func init() {
	actorSystemInstance.actorNames = make(map[string]actorRefInterface)
	actorSystemInstance.actors = make(map[actorRefInterface]actorInterface)
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

type actorSystem struct {
	actorNames map[string]actorRefInterface
	actors     map[actorRefInterface]actorInterface
}

func (system *actorSystem) RegisterActor(name string, actor actorInterface) error {
	_, exists := system.actorNames[name]
	if exists {
		return fmt.Errorf("Actor %v already registered")
	}

	actor.setName(name)
	actorRef := &ActorRef{name}
	system.actorNames[name] = *actorRef
	system.actors[*actorRef] = actor

	actor.setMailbox(make(chan Message))
	go receive(actor)

	return nil
}

func (system *actorSystem) Actor(name string) (actorRefInterface, error) {
	ref, exists := system.actorNames[name]
	if !exists {
		return nil, fmt.Errorf("Actor %v not registered", name)
	}

	return ref, nil
}

func (system *actorSystem) actor(actorRef actorRefInterface) (actorInterface, error) {
	ref, exists := system.actors[actorRef]

	if !exists {
		return nil, fmt.Errorf("Actor %v not registered", actorRef)
	}

	return ref, nil
}
