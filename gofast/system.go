package gofast

import (
	"fmt"
	"Gofast/gofast/util"
)

var actorSystemInstance actorSystem
var distributedActorConfiguration map[string]DistributedInterface

func InitLocalActorSystem() {
	actorSystemInstance = actorSystem{}
	actorSystemInstance.actors = make(map[string]actorAssociation)
	distributedActorConfiguration = make(map[string]DistributedInterface)
}

func InitDistributedActorSystem(remote ... DistributedInterface) error {
	InitLocalActorSystem()

	for _, v := range remote {
		alias := v.ConnectionAlias()
		_, exists := distributedActorConfiguration[alias]
		if !exists {
			err := v.Connection()
			if err != nil {
				util.LogFatal("Failed to create the distributed actor system: %v", err)
				return err
			}
			distributedActorConfiguration[v.ConnectionAlias()] = v
		}
	}

	//etcd.InitConfiguration(endpoints...)

	return nil
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

type DistributedInterface interface {
	Configure(string, string)
	Connection() error
	Send(string, []byte)
	Receive(string)
	Close()
	ConnectionAlias() string
}

type actorAssociation struct {
	actorRef ActorRefInterface
	actor    actorInterface
	options  OptionsInterface
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
		actorAssociation{actorRef, actor, options}

	go receive(actor, options)

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

func (system *actorSystem) DistributedConfiguration(connectionAlias string) DistributedInterface {
	v, _ := distributedActorConfiguration[connectionAlias]
	return v
}

func (system *actorSystem) Invoke(message Message) error {
	actorAssociation, err := system.actor(message.Self.Name())

	if err != nil {
		util.LogError("Actor %v not registered")
		return err
	}

	if message.messageType == poisonPill {
		util.LogInfo("Actor %v has received a poison pill", actorAssociation.actor.Name())
		actorAssociation.actor.Close()
		return nil
	}

	f, exists := actorAssociation.actor.reactions()[message.messageType]
	if exists {
		f(message)
	}

	return nil
}
