package gosiris

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"time"
)

var root *Actor
var actorSystemInstance actorSystem
var actorSystemStarted bool
var registry registryInterface

type SystemOptions struct {
	ActorSystemName string
	ZipkinOptions   ZipkinOptions
	RegistryUrl     string
}

func init() {
	root = &Actor{}
	root.name = "root"
	actorSystemStarted = false
}

func InitActorSystem(options SystemOptions) error {
	if actorSystemStarted {
		ErrorLogger.Printf("Actor system already started")
		return fmt.Errorf("actor system already started")
	}

	actorSystemInstance = actorSystem{}
	actorSystemInstance.actors = make(map[string]actorAssociation)
	actorSystemStarted = true

	if options.RegistryUrl != "" {
		initDistributedActorSystem(options.RegistryUrl)
	}

	if options.ZipkinOptions.Url != "" {
		initZipkinSystem(options.ActorSystemName, options.ZipkinOptions)
	}

	return nil
}

func initDistributedActorSystem(url string) error {
	//TODO Manage the implementation dynamically
	registry = &etcdClient{}
	err := registry.Configure(url)
	if err != nil {
		actorSystemStarted = false
		FatalLogger.Fatalf("Failed to configure access to the remote repository: %v", err)
	}

	conf, err := registry.ParseConfiguration()
	if err != nil {
		actorSystemStarted = false
		FatalLogger.Fatalf("Failed to parse the configuration: %v", err)
	}

	InitRemoteConnections(conf)
	ActorSystem().addRemoteActors(conf)

	return nil
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

func CloseActorSystem() error {
	if !actorSystemStarted {
		ErrorLogger.Printf("Actor system not started")
		return fmt.Errorf("actor system not started")
	}

	if registry != nil {
		registry.Close()
	}
	InfoLogger.Printf("Actor system closed")

	actorSystemStarted = false
	closeZipkinSystem()

	return nil
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
	if name == root.name {
		ErrorLogger.Printf("Register an actor whose name is %v is not allowed", name)
		return fmt.Errorf("register an actor whose name is %v is not allowed", name)
	}

	InfoLogger.Printf("Registering new actor %v", name)
	return system.SpawnActor(RootActor(), name, actor, options)
}

func (system *actorSystem) SpawnActor(parent actorInterface, name string, actor actorInterface, options OptionsInterface) error {
	InfoLogger.Printf("Spawning new actor %v", name)

	_, exists := system.actors[name]
	if exists {
		InfoLogger.Printf("Actor %v already registered", name)
	}

	if options == nil {
		options = &ActorOptions{}
		options.SetRemote(false)
	}

	if options.BufferSize() == 0 {
		options.SetBufferSize(defaultBufferSize)
	}

	actor.setName(name)
	actor.setParent(parent)
	options.setParent(parent.Name())
	if !options.Remote() {
		actor.setDataChan(make(chan Context, options.BufferSize()))
		actor.setCloseChan(make(chan interface{}))
	} else {
		registry.RegisterActor(name, options)
		go registry.Watch(system.onActorCreatedFromRegistry, system.onActorRemovedFromRegistry)
		AddConnection(name, options)
	}

	actorRef := newActorRef(name)

	system.actors[name] =
		actorAssociation{actorRef, actor, options}

	go receive(actor, options)

	if options.DefaultWatcher() != 0 {
		//Configure a default watcher
		t := time.NewTicker(options.DefaultWatcher())
		go func(t *time.Ticker, actorRef ActorRefInterface) {
			for {
				select {
				case <-t.C:
					p, err := system.actor(parent.Name())
					if err != nil {
						ErrorLogger.Printf("Parent of actor %v not found", name)
					}
					dispatch(p.actor.getDataChan(), GosirisMsgHeartbeatRequest, nil, actorRef, p.actorRef, new(ActorOptions), nil)
				}
			}
		}(t, actorRef)
	}

	return nil
}

func (system *actorSystem) onActorCreatedFromRegistry(name string, options *ActorOptions) {
	actorRef := newActorRef(name)
	actor := Actor{}
	actor.name = name

	system.actors[name] =
		actorAssociation{actorRef, nil, options}

	AddConnection(name, options)

	InfoLogger.Printf("Actor %v added to the local system", name)
}

func (system *actorSystem) onActorRemovedFromRegistry(name string) {
	system.removeRemoteActor(name)

	InfoLogger.Printf("Actor %v removed from the local system", name)
}

func (system *actorSystem) removeRemoteActor(name string) {
	InfoLogger.Printf("Removing remote actor %v", name)

	v, err := system.actor(name)

	if err != nil {
		InfoLogger.Printf("Actor %v not registered", name)
		return
	}

	if v.options.Remote() {
		DeleteRemoteActorConnection(name)
	}

	delete(system.actors, name)

	InfoLogger.Printf("Remote actor %v removed", name)
}

func (system *actorSystem) closeLocalActor(name string) {
	InfoLogger.Printf("Closing local actor %v", name)

	v, err := system.actor(name)

	if err != nil {
		ErrorLogger.Printf("Unable to close actor %v", name)
		return
	}

	//If the actor has a parent we send him a message
	if v.actor.Parent() != nil {
		parentName := v.actor.Parent().Name()
		if parentName != root.name {
			p, err := system.actor(parentName)
			if err != nil {
				ErrorLogger.Printf("Parent %v not registered", parentName)
			}
			p.actorRef.Tell(EmptyContext, GosirisMsgChildClosed, name, v.actorRef)
		}
	}

	//m := v.actor.getDataChan()
	//if m != nil {
	//	close(m)
	//}

	c := v.actor.getCloseChan()
	if c != nil {
		c <- 0
	}

	if registry != nil {
		registry.UnregisterActor(name)
	}

	delete(system.actors, name)

	InfoLogger.Printf("%v unregistered from the actor system", name)
}

func (system *actorSystem) actor(name string) (actorAssociation, error) {
	ref, exists := system.actors[name]
	if !exists {
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

	InfoLogger.Printf("Actors configuration: %v", system.actors)
}

func (system *actorSystem) ActorOf(name string) (ActorRefInterface, error) {
	actorAssociation, err := system.actor(name)

	if err != nil {
		return nil, err
	}

	return actorAssociation.actorRef, err
}

func (system *actorSystem) Invoke(message Context) error {
	if message.Self == nil {
		return nil
	}

	InfoLogger.Printf("Invoking %v", message)

	actorAssociation, err := system.actor(message.Self.Name())

	if err != nil {
		ErrorLogger.Printf("Invoke error: actor %v not registered", actorAssociation.actorRef.Name())
		return err
	}

	if message.MessageType == GosirisMsgPoisonPill {
		InfoLogger.Printf("Actor %v has received a poison pill", actorAssociation.actor.Name())

		if actorAssociation.options.Autoclose() {
			InfoLogger.Printf("Autoclose actor %v", actorAssociation.actor.Name())
			actorAssociation.actor.Close()
			return nil
		}
	} else if message.MessageType == GosirisMsgHeartbeatRequest {
		InfoLogger.Printf("Actor %v has received a heartbeat request", actorAssociation.actor.Name())

		message.Sender.Tell(EmptyContext, GosirisMsgHeartbeatReply, nil, message.Self)
	}

	if actorAssociation.actor.reactions() != nil {
		f, exists :=
			actorAssociation.actor.reactions()[message.MessageType]
		if exists {
			var span opentracing.Span
			if message.carrier != nil {
				ctx, _ := extract(message.carrier)
				span = tracer.StartSpan("operation", opentracing.ChildOf(ctx))
				InfoLogger.Printf("Starting child span")
				message.span = span
			}
			f(message)
			if message.carrier != nil {
				span.Finish()
			}
		}
	}

	return nil
}

func (system *actorSystem) Stop(stop chan struct{}) {
	if stop != nil {
		stop <- struct{}{}
	}
}
