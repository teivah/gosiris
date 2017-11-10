[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard)

gosiris is an [actor](https://en.wikipedia.org/wiki/Actor_model) framework for Golang.

# Features

* Manage a hierarchy of actors (each actor has its own: state, behavior, mailbox, child actors)
* Deploy remote actors accessible though an AMQP broker or Kafka
* Automated registration and runtime discoverability using etcd registry
* Zipkin integration 
* Built-in patterns (become/unbecome, send, forward, repeat, child supervision)

# Examples

## Hello world

```go
package main

import (
	"gosiris/gosiris"
)

func main() {
	//Init a local actor system
	gosiris.InitActorSystem(gosiris.SystemOptions{
		ActorSystemName: "ActorSystem",
	})

	//Create an actor
	parentActor := gosiris.Actor{}
	//Close an actor
	defer parentActor.Close()

	//Create an actor
	childActor := gosiris.Actor{}
	//Close an actor
	defer childActor.Close()
	//Register a reaction to event types ("message" in this case)
	childActor.React("message", func(context gosiris.Context) {
		context.Self.LogInfo(context, "Received %v\n", context.Data)
	})

	//Register an actor to the system
	gosiris.ActorSystem().RegisterActor("parentActor", &parentActor, nil)
	//Register an actor by spawning it
	gosiris.ActorSystem().SpawnActor(&parentActor, "childActor", &childActor, nil)

	//Retrieve actor references
	parentActorRef, _ := gosiris.ActorSystem().ActorOf("parentActor")
	childActorRef, _ := gosiris.ActorSystem().ActorOf("childActor")

	//Send a message from one actor to another (from parentActor to childActor)
	childActorRef.Tell(gosiris.EmptyContext, "message", "Hi! How are you?", parentActorRef)
}
```

```
INFO: [childActor] 1988/01/08 01:00:00 Received Hi! How are you?
```

## Distributed actor system example

In the following example, **in less than 30 effective lines of code**, we will see how to create a distributed actor system implementing a request/reply interaction. An actor will be triggered by AMQP messages while another one will be triggered by Kafka events. Each actor will register itself in an etcd instance and will discover the other actor at runtime. Last but not least, gosiris will also manage the Zipkin integration by automatically managing the spans and forwarding the logs.

```go
package main

import (
	"gosiris/gosiris"
	"time"
)

func main() {
	//Configure a distributed actor system with an etcd registry and a Zipkin integration
	gosiris.InitActorSystem(gosiris.SystemOptions{
		ActorSystemName: "ActorSystem",
		RegistryUrl:     "http://etcd:2379",
		ZipkinOptions: gosiris.ZipkinOptions{
			Url:      "http://zipkin:9411/api/v1/spans",
			Debug:    true,
			HostPort: "0.0.0.0",
			SameSpan: true,
		},
	})
	//Defer the actor system closure
	defer gosiris.CloseActorSystem()

	//Configure actor1
	actor1 := new(gosiris.Actor).React("reply", func(context gosiris.Context) {
		//Because Zipkin is enabled, the log will be also sent to the Zipkin server
		context.Self.LogInfo(context, "Received: %v", context.Data)

	})
	//Defer actor1 closure
	defer actor1.Close()
	//Register a remote actor accessible through AMQP
	gosiris.ActorSystem().RegisterActor("actor1", actor1, new(gosiris.ActorOptions).SetRemote(true).SetRemoteType(gosiris.Amqp).SetUrl("amqp://guest:guest@amqp:5672/").SetDestination("actor1"))

	//Configure actor2
	actor2 := new(gosiris.Actor).React("context", func(context gosiris.Context) {
		//Because Zipkin is enabled, the log will be also sent to the Zipkin server
		context.Self.LogInfo(context, "Received: %v", context.Data)
		context.Sender.Tell(context, "reply", "hello back", context.Self)
	})
	//Defer actor2 closure
	defer actor2.Close()
	//Register a remote actor accessible through Kafka
	gosiris.ActorSystem().SpawnActor(actor1, "actor2", actor2, new(gosiris.ActorOptions).SetRemote(true).SetRemoteType(gosiris.Kafka).SetUrl("kafka:9092").SetDestination("actor2"))

	//Retrieve the actor references
	actor1Ref, _ := gosiris.ActorSystem().ActorOf("actor1")
	actor2Ref, _ := gosiris.ActorSystem().ActorOf("actor2")

	//Send a message to the kafkaRef
	actor2Ref.Tell(gosiris.EmptyContext, "context", "hello", actor1Ref)

	time.Sleep(250 * time.Millisecond)
}

```

```
INFO: [actor2] 2017/11/11 00:38:24 Received: hello
INFO: [actor1] 2017/11/11 00:38:24 Received: hello back
```

## More Examples

See the examples in [actor_test.go](gosiris/actor_test.go).

# Environment

To setup the complete gosiris environment:
* etcd:

```bash
docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 quay.io/coreos/etcd:v2.3.8 -name etcd0 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -initial-advertise-peer-urls http://${HostIP}:2380 -listen-peer-urls http://0.0.0.0:2380 -initial-cluster-token etcd-cluster-1 -initial-cluster etcd0=http://${HostIP}:2380 -initial-cluster-state new
```

* An AMQP broker (e.g. RabbitMQ):
 
```bash
docker run -d --hostname rabbit --name rabbit -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 15672:15672 rabbitmq
docker exec rabbit rabbitmq-plugins enable rabbitmq_management
```

The last command is not mandatory but it allows to expose a web UI on the port 15672.

* A Kafka broker:

[https://teivah.io/blog/running-kafka-1-0-in-docker/](https://teivah.io/blog/running-kafka-1-0-in-docker/)

* A Zipkin server:

docker run --name zipkin -d -p 9411:9411 openzipkin/zipkin

Meanwhile, the gosiris tests are using several hostnames you need to configure: _etcd_, _amqp_, _zipkin_, and _kafka_.

# Troubleshooting

You may experience errors during the tests like the following:
```
r.EncodeArrayStart undefined (type codec.encDriver has no field or method EncodeArrayStart)
```

This is a known issue with the etcd client used. The manual workaround (for the time being) is to delete manually the file _keys.generated.go_ generated in /vendor. 

# Contributing

* Open an issue if you want a new feature or if you spotted a bug
* Feel free to propose pull requests

Any contribution is more than welcome! In the meantime, if we want to discuss gosiris you can contact me [@teivah](https://twitter.com/teivah).