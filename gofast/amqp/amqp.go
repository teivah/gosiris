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

func Close() {
	channel.Close()
	connection.Close()
}

func Publish() {
	//q, err := ch.QueueDeclare(
	//	"hello", // name
	//	false,   // durable
	//	false,   // delete when unused
	//	false,   // exclusive
	//	false,   // no-wait
	//	nil,     // arguments
	//)
	//if err != nil {
	//	util.errorLogger.Printf()
	//}
	//
	//body := "hello"
	//err = ch.Publish(
	//	"",     // exchange
	//	q.Name, // routing key
	//	false,  // mandatory
	//	false,  // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	//failOnError(err, "Failed to publish a message")
}
