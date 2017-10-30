gopera is a simple library to bring the actor model on top of Golang.

# Features
* Send message from one actor to another using the **mailbox** principle
* **Forward** message to maintain the original sender
* **Hierarchy** concept between the different actors
* **Become/unbecome** principle to modify at runtime the behavior of an actor
* Capacity to gracefully ask an actor to **stop** its execution
* **Distributed actor system** across the network using an **AMQP broker**
* Actors **discoverability** using etcd 

# Hello world
The gofast hello world is the following:

```go
package main

import (
	"github.com/teivah/gopera/gopera"
)

func main() {
	//Init a local actor system
	gopera.InitLocalActorSystem()

	//Create an actor
	parentActor := gopera.Actor{}
	//Close an actor
	defer parentActor.Close()

	//Create an actor
	childActor := gopera.Actor{}
	//Close an actor
	defer childActor.Close()
	//Register a reaction to event types ("message" in this case)
	childActor.React("message", func(message gopera.Message) {
		message.Self.LogInfo("Received %v\n", message.Data)
	})

	//Register an actor to the system
	gopera.ActorSystem().RegisterActor("parentActor", &parentActor, nil)
	//Register an actor by spawning it
	gopera.ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	//Retrieve actor references
	parentActorRef, _ := gopera.ActorSystem().ActorOf("parentActor")
	childActorRef, _ := gopera.ActorSystem().ActorOf("childActor")

	//Send a message from one actor to another (from parentActor to childActor)
	childActorRef.Send("message", "Hi! How are you?", parentActorRef)
}
```

```
[childActor] Received Hi! How are you?
```

# Contributing

* Open an issue if you need a new feature or if you spotted a bug
* Feel free to propose pull requests

Any contribution is welcome!