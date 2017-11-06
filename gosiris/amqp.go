package gosiris

import (
	"github.com/streadway/amqp"
	"encoding/json"
)

type AmqpConnection struct {
	url        string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (amqpConnection *AmqpConnection) Configure(url string) {
	amqpConnection.url = url
}

func (amqpConnection *AmqpConnection) Connection() error {
	c, err := amqp.Dial(amqpConnection.url)
	if err != nil {
		ErrorLogger.Printf("Failed to connect to the AMQP server %v", amqpConnection.url)
		return err
	}
	amqpConnection.connection = c

	ch, err := c.Channel()
	if err != nil {
		ErrorLogger.Printf("Failed to open an AMQP channel on the server %v", amqpConnection.url)
		return err
	}
	amqpConnection.channel = ch

	InfoLogger.Printf("Connected to %v", amqpConnection.url)

	return nil
}

func (amqpConnection *AmqpConnection) Receive(queueName string) {
	q, err := amqpConnection.channel.QueueDeclare(
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

	msgs, err := amqpConnection.channel.Consume(
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

func (amqpConnection *AmqpConnection) Close() {
	amqpConnection.channel.Close()
	amqpConnection.connection.Close()
}

func (amqpConnection *AmqpConnection) Send(destination string, data []byte) {
	json := string(data)
	InfoLogger.Printf("Sending %v", json)

	q, err := amqpConnection.channel.QueueDeclare(
		destination, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		ErrorLogger.Printf("Error while declaring queue %v: %v", destination, err)
	}

	body := data
	err = amqpConnection.channel.Publish(
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
	}
}
