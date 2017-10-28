package actor

import "fmt"

var systemInstance system = system{}

func System() system {
	return systemInstance
}


type system struct {
	conf map[string]actorInterface
}

type SystemActor interface {
	RegisterActor(actorInterface) error
	Configuration() map[string]actorInterface
}

func (system *system) RegisterActor(actor actorInterface) error {
	if system.conf == nil {
		system.conf = make(map[string]actorInterface)
	}

	if system.conf[actor.Name()] != nil {
		return fmt.Errorf("Actor %v already registered")
	}

	system.conf[actor.Name()] = actor
	return nil
}

func (system *system) Configuration() map[string]actorInterface {
	return system.conf
}