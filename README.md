# gofast

gofast is a simple library to bring the actor model on top of Golang.

## Features
* Send message from one actor to another using the **mailbox** principle
* **Forward** message to maintain the original sender
* **Hierarchy** concept between the different actors
* **Become/unbecome** principle to modify at runtime the behavior of an actor
* Capacity to gracefully ask an actor to **stop** its execution
* **Distributed actor system** across the network using an **AMQP broker**

## Hello world
The gofast hello world is the following:

```go
//Create a simple parent actor
parentActor := Actor{}
//Close the actor
defer parentActor.Close()

//Register the parent actor
ActorSystem().RegisterActor("parentActor", &parentActor, nil)

//Create a simple child actor
childActor := Actor{}
//Close the actor
defer childActor.Close()

//Register the reactions to event types (here a reaction to message)
childActor.React("message", func(message Message) {
    message.Self.LogInfo("Received %v\n", message.Data)
})

//Register the child actor
ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

//Retrieve the parent and child actor reference
parentActorRef, _ := ActorSystem().ActorOf("parentActor")
childActorRef, _ := ActorSystem().ActorOf("childActor")

//Send a message from the parent to the child actor
childActorRef.Send("message", "Hi! How are you?", parentActorRef)
```

```
[childActor] Received Hi! How are you?
```

## Participation

If you want to participate, feel free to contact me contact me [@teivah](https://twitter.com/teivah)