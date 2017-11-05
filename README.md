gopera is an [actor](https://en.wikipedia.org/wiki/Actor_model) framework for Golang.

# Principles
gopera is based on three principles: **Configure**, **Discover**, **React**
* Configure the actor behaviour depending on event types
* Discover the other actors dynamically registered in a remote registry
* React on events sent by actors

In a nutshell, gopera allows to:
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
INFO: [childActor] 1988/01/08 01:00:00 Received Hi! How are you?
```

# Features

## Message forwarding

```go
//Create an actor
actor := gopera.Actor{}
//Configure it to react on someMessage
actor.React("someMessage", func(message gopera.Message) {
    //Forward the message to fooActor and barActor
    message.Self.Forward(message, "fooActor", "barActor")
})
```

## Stateful actor

```go
//Create a structure extending the standard gopera.Actor one
type StatefulActor struct {
	gopera.Actor
	someValue int
}

//...

//Create an actor
actor := StatefulActor{}
//Configure it to react on someMessage
actor.React("someMessage", func(message gopera.Message) {
    //Modify the actor internal state
    if message.Data == 0 {
        actor.someValue = 1
    } else {
        actor.someValue = 0
    }
})
```

## Become/unbecome

```go
//Bar behavior
bar := func(message gopera.Message) {
    if message.Data == "foo" {
        //Unbecome to foo
        message.Self.Unbecome(message.MessageType)
    }
}

//Foo behavior
foo := func(message gopera.Message) {
    if message.Data == "bar" {
        //Become bar
        message.Self.Become(message.MessageType, bar)
    }
}

//Create an actor
actor := gopera.Actor{}
//Configure it to react on someMessage
actor.React("someMessage", foo)
```

```go
```

# Contributing

* Open an issue if you want a new feature or if you spotted a bug
* Feel free to propose pull requests

Any contribution is welcome!