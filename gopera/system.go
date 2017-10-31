package gopera

import (
	"fmt"
	"gopera/gopera/util"
)

var actorSystemInstance actorSystem
var registry registryInterface

func InitLocalActorSystem() {
	actorSystemInstance = actorSystem{}
	actorSystemInstance.actors = make(map[string]actorAssociation)
}

func InitDistributedActorSystem(url ...string) error {
	InitLocalActorSystem()

	//TODO Manage the implementation dynamically
	registry = &etcdClient{}
	err := registry.Configure(url...)
	if err != nil {
		util.LogFatal("Failed to configure access to the remote repository: %v", err)
	}

	conf, err := registry.ParseConfiguration()
	if err != nil {
		util.LogFatal("Failed to parse the configuration: %v", err)
	}

	InitRemoteConnections(conf)
	ActorSystem().addRemoteActors(conf)

	return nil
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
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
		util.LogInfo("Actor %v already registered", name)
	}

	if options == nil {
		options = &ActorOptions{}
		options.SetRemote(false)
	}

	actor.setName(name)
	actor.setParent(parent)
	if !options.Remote() {
		actor.setMailbox(make(chan Message))
	} else {
		registry.RegisterActor(name, options)
		go registry.Watch(system.cbCreate, system.cbDelete)
		AddConnection(name, options)
	}

	actorRef := newActorRef(name)

	system.actors[name] =
		actorAssociation{actorRef, actor, options}

	go receive(actor, options)

	return nil
}

func (system *actorSystem) cbCreate(name string, options *ActorOptions) {
	actorRef := newActorRef(name)
	actor := Actor{}
	actor.name = name

	system.actors[name] =
		actorAssociation{actorRef, &actor, options}

	AddConnection(name, options)

	util.LogInfo("Actor %v added to the local system", name)
}

func (system *actorSystem) cbDelete(name string) {
	delete(system.actors, name)

	util.LogInfo("Actor %v removed from the local system", name)
}

func (system *actorSystem) unregisterActor(name string) {
	_, err := system.ActorOf(name)
	if err != nil {
		util.LogError("Actor %v not registered", name)
		return
	}

	v, err := system.actor(name)
	if v.options.Remote() {
		DeleteConnection(name)
	}

	registry.UnregisterActor(name)

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

func (system *actorSystem) addRemoteActors(configuration map[string]OptionsInterface) {
	for k, v := range configuration {
		actor := Actor{}
		actor.setName(k)
		actorRef := newActorRef(k)

		system.actors[k] = actorAssociation{actorRef, &actor, v}
	}

	util.LogInfo("Actors configuration: %v", system.actors)
}

func (system *actorSystem) ActorOf(name string) (ActorRefInterface, error) {
	actorAssociation, err := system.actor(name)

	if err != nil {
		return nil, err
	}

	return actorAssociation.actorRef, err
}

func (system *actorSystem) Invoke(message Message) error {
	actorAssociation, err := system.actor(message.Self.Name())

	if err != nil {
		util.LogError("Actor %v not registered")
		return err
	}

	if message.messageType == poisonPill {
		util.LogInfo("Actor %v has received a poison pill", actorAssociation.actor.Name())

		if actorAssociation.options.Autoclose() {
			util.LogInfo("Autoclose actor %v", actorAssociation.actor.Name())
			actorAssociation.actor.Close()
		}
		return nil
	}

	f, exists := actorAssociation.actor.reactions()[message.messageType]
	if exists {
		f(message)
	}

	return nil
}
