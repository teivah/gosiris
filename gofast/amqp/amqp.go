package amqp

import (
	"fmt"
	"Gofast/gofast/util"
	"github.com/streadway/amqp"
)

type RemoteAmqp struct {
	connectionAlias string
	url             string
	connection      *amqp.Connection
	channel         *amqp.Channel
}

func (remoteAmqp *RemoteAmqp) Configure(connectionAlias string, url string) {
	remoteAmqp.connectionAlias = connectionAlias
	remoteAmqp.url = url
}

func (remoteAmqp *RemoteAmqp) Connection() error {
	c, err := amqp.Dial(remoteAmqp.url)
	if err != nil {
		util.LogError("Failed to connect to the AMQP server %v", remoteAmqp.url)
		return err
	}
	remoteAmqp.connection = c

	ch, err := c.Channel()
	if err != nil {
		util.LogError("Failed to open an AMQP channel on the server %v", remoteAmqp.url)
		return err
	}
	remoteAmqp.channel = ch

	return nil
}

func (remoteAmqp *RemoteAmqp) Receive(queueName string) {
	q, err := remoteAmqp.channel.QueueDeclare(
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

	msgs, err := remoteAmqp.channel.Consume(
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

func (remoteAmqp *RemoteAmqp) Close() {
	remoteAmqp.channel.Close()
	remoteAmqp.connection.Close()
}

func (remoteAmqp *RemoteAmqp) Send(destination string) {
	q, err := remoteAmqp.channel.QueueDeclare(
		destination, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		util.LogError("Error while declaring queue %v: %v", destination, err)
	}

	body := "hello"
	err = remoteAmqp.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		util.LogError("Error while publishing a message to queue %v: %v", destination, err)
	}
}

func (remoteAmqp *RemoteAmqp) ConnectionAlias() string {
	return remoteAmqp.connectionAlias
}
