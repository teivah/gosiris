package gofast

import (
	"fmt"
)

var actorSystemInstance actorSystem = actorSystem{}

func init() {
	actorSystemInstance.actorNames = make(map[string]ActorRefInterface)
	actorSystemInstance.actors = make(map[ActorRefInterface]actorInterface)
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

type actorSystem struct {
	actorNames map[string]ActorRefInterface
	actors     map[ActorRefInterface]actorInterface
}

func (system *actorSystem) RegisterActor(name string, actor actorInterface) error {
	return system.SpawnActor(RootActor(), name, actor)
}

func (system *actorSystem) SpawnActor(parent actorInterface, name string, actor actorInterface) error {
	_, exists := system.actorNames[name]
	if exists {
		return fmt.Errorf("actor %v already registered", name)
	}

	actor.setName(name)
	actor.setParent(parent)
	actor.setMailbox(make(chan Message))

	actorRef :=
		ActorRef{name}
	system.actorNames[name] = actorRef
	system.actors[actorRef] = actor

	go receive(actor)

	return nil
}

func (system *actorSystem) unregisterActor(name string) {
	actorRef, exists := system.actorNames[name]
	if !exists {
		return
	}

	delete(system.actorNames, name)
	delete(system.actors, actorRef)
}

func (system *actorSystem) Actor(name string) (ActorRefInterface, error) {
	ref, exists := system.actorNames[name]
	if !exists {
		return nil, fmt.Errorf("actor %v not registered", name)
	}

	return ref, nil
}

func (system *actorSystem) actor(actorRef ActorRefInterface) (actorInterface, error) {
	ref, exists := system.actors[actorRef]

	if !exists {
		return nil, fmt.Errorf("actor %v not registered", actorRef)
	}

	return ref, nil
}
