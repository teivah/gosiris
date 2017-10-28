# Gofast

Gofast is a simple library to bring the actor model on top of Golang.

Just like [Akka](https://akka.io/) for example, Gofast implements a hierarchy concept between the actors to improve the supervision.

A hello world example would be the following:

```go
//Create a simple parent actor
parentActor := Actor{}
//Register the parent actor
ActorSystem().RegisterActor("parentActor", &parentActor, Root())

//Create a simple child actor
childActor := Actor{}
//Register the reactions to event types (here a reaction to message)
childActor.React("message", func(message Message) {
	childActor.Printf("Received %v\n", message.Data)
})
//Register the child actor
ActorSystem().RegisterActor("childActor", &childActor, &parentActor)

//Retrieve the parent and child actor reference
parentActorRef, _ := ActorSystem().Actor("parentActor")
childActorRef, _ := ActorSystem().Actor("childActor")

//Send a message from the parent to the child actor
childActorRef.Tell("message", "Hi! How are you?", parentActorRef)
```

```
[childActor] Received [Hi! How are you?]
```

We would like to add many other features like the capacity to distribute actors across the network etc. If you want to participate, contact me [@teivah](https://twitter.com/teivah)