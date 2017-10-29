package amqp

import (
	"github.com/streadway/amqp"
	"fmt"
	"Gofast/gofast/util"
)

var connection *amqp.Connection
var channel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		util.LogFatal("%v: %v", msg, err)
		panic(fmt.Sprintf("%v: %v", msg, err))
	}
}

func InitConfiguration(endpoint string) {
	c, err := amqp.Dial(endpoint)
	failOnError(err, "Failed to connect to RabbitMQ")
	connection = c

	ch, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	channel = ch
}

func AddConsumer(queueName string) {
	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		util.LogError("Error while declaring queue %v: %v", queueName, err)
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s", d.Body)
		}
	}()
}

func Close() {
	channel.Close()
	connection.Close()
}

func Publish(queueName string) {
	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		util.LogError("Error while declaring queue %v: %v", queueName, err)
	}

	body := "hello"
	err = channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		util.LogError("Error while publishing a message to queue %v: %v", queueName, err)
	}
}
