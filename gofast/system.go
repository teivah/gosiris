package gofast

import (
	"fmt"
	"Gofast/gofast/util"
)

var actorSystemInstance actorSystem
var remoteConfiguration map[string]RemoteInterface

func InitLocalActorSystem() {
	actorSystemInstance = actorSystem{}
	actorSystemInstance.actors = make(map[string]actorAssociation)
	remoteConfiguration = make(map[string]RemoteInterface)
}

func InitDistributedActorSystem(remote ... RemoteInterface) error {
	InitLocalActorSystem()

	for _, v := range remote {
		alias := v.ConnectionAlias()
		_, exists := remoteConfiguration[alias]
		if !exists {
			err := v.Connection()
			if err != nil {
				util.LogFatal("Failed to create the distributed actor system: %v", err)
				return err
			}
		}
	}

	//etcd.InitConfiguration(endpoints...)

	return nil
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

type RemoteInterface interface {
	Configure(string, string)
	Connection() error
	Send(string)
	Receive(string)
	Close()
	ConnectionAlias() string
}

type actorAssociation struct {
	actorRef ActorRefInterface
	actor    actorInterface
}

type actorSystem struct {
	actors map[string]actorAssociation
}

func (system *actorSystem) RegisterActor(name string, actor actorInterface, options OptionsInterface) error {
	util.LogInfo("Registering new actor %v", name)
	return system.SpawnActor(RootActor(), name, actor, options)
}

func (system *actorSystem) SpawnActor(parent actorInterface, name string, actor actorInterface, options OptionsInterface) error {
	util.LogInfo("Spawning new actor %v", name)

	_, exists := system.actors[name]
	if exists {
		util.LogInfo("actor %v already registered", name)
		return fmt.Errorf("actor %v already registered", name)
	}

	actor.setName(name)
	actor.setParent(parent)
	actor.setMailbox(make(chan Message))

	actorRef := newActorRef(name)

	system.actors[name] =
		actorAssociation{actorRef, actor}

	go receive(actor)

	return nil
}

func (system *actorSystem) unregisterActor(name string) {
	_, err := system.ActorOf(name)
	if err != nil {
		util.LogError("Actor %v not registered", name)
		return
	}

	delete(system.actors, name)

	util.LogInfo("%v unregistered from the actor system", name)
}

func (system *actorSystem) actor(name string) (actorAssociation, error) {
	ref, exists := system.actors[name]
	if !exists {
		util.LogError("actor %v not registered", name)
		return actorAssociation{}, fmt.Errorf("actor %v not registered", name)
	}

	return ref, nil
}

func (system *actorSystem) ActorOf(name string) (ActorRefInterface, error) {
	actorAssociation, err := system.actor(name)

	if err != nil {
		return nil, err
	}

	return actorAssociation.actorRef, err
}
