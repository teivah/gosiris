gopera is an [actor](https://en.wikipedia.org/wiki/Actor_model) framework for Golang.

# Principles
gopera is based on three principles: **Configure**, **Discover**, **React**
* Configure the actor behaviour depending on event types
* Discover the other actors dynamically registered in a remote registry
* React on events sent by actors

# Detailed features
* Send messages from one actor to another using the **mailbox** principle.
* An actor can be either **local** (triggered through a Go channel) or **remote** (triggered through an **AMQP broker**)
* **Hierarchy** dependencies between the different actors
* **Forward** message to maintain the original sender
* **Become/unbecome** capability to modify at runtime the behavior of an actor
* Capacity to gracefully ask an actor to **stop** its execution
* Automatic registration and **discoverability** of the actors using a registry (**etcd**)

# Hello world

```go
package main

import (
	"github.com/teivah/gopera/gopera"
)

func main() {
	//Init a local actor system
    gopera.InitLocalActorSystem()
    //Defer the actor system closure
    defer gopera.CloseActorSystem()

    //Create a simple parent actor
    parentActor := gopera.Actor{}
    //Defer the actor closure
    defer parentActor.Close()

    //Register the parent actor
    gopera.ActorSystem().RegisterActor("parentActor", &parentActor, nil)

    //Create a simple child actor
    childActor := gopera.Actor{}
    //Defer the actor system closure
    defer childActor.Close()

    //Register the reactions to event types (here a reaction to message)
    childActor.React("message", func(message gopera.Message) {
        message.Self.LogInfo("Received %v\n", message.Data)
    })

    //Register the child actor
    gopera.ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

    //Retrieve the parent and child actor reference
    parentActorRef, _ := gopera.ActorSystem().ActorOf("parentActor")
    childActorRef, _ := gopera.ActorSystem().ActorOf("childActor")

    //Tell a message from the parent to the child actor
    childActorRef.Tell("message", "Hi! How are you?", parentActorRef)
}
```

```
INFO: [childActor] 1988/01/08 01:00:00 Received Hi! How are you?
```

# Contributing

* Open an issue if you want a new feature or if you spotted a bug
* Feel free to propose pull requests

Any contribution is welcome!