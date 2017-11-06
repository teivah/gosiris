package gosiris

import (
	"github.com/streadway/amqp"
	"encoding/json"
)

var Amqp = "amqp"

type amqpTransport struct {
	url        string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (a *amqpTransport) Configure(url string, options map[string]string) {
	a.url = url
}

func (a *amqpTransport) Connection() error {
	c, err := amqp.Dial(a.url)
	if err != nil {
		ErrorLogger.Printf("Failed to connect to the AMQP server %v", a.url)
		return err
	}
	a.connection = c

	ch, err := c.Channel()
	if err != nil {
		ErrorLogger.Printf("Failed to open an AMQP channel on the server %v", a.url)
		return err
	}
	a.channel = ch

	InfoLogger.Printf("Connected to %v", a.url)

	return nil
}

func (a *amqpTransport) Receive(queueName string) {
	q, err := a.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		ErrorLogger.Printf("Error while declaring queue %v: %v", queueName, err)
	}

	msgs, err := a.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	for d := range msgs {
		msg := Message{}
		json.Unmarshal(d.Body, &msg)

		ActorSystem().Invoke(msg)
	}
}

func (a *amqpTransport) Close() {
	a.channel.Close()
	a.connection.Close()
}

func (a *amqpTransport) Send(destination string, data []byte) error {
	json := string(data)
	InfoLogger.Printf("Sending %v", json)

	q, err := a.channel.QueueDeclare(
		destination, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		ErrorLogger.Printf("Error while declaring queue %v: %v", destination, err)
		return err
	}

	body := data
	err = a.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		ErrorLogger.Printf("Error while publishing a message to queue %v: %v", destination, err)
		return err
	}

	return nil
}
