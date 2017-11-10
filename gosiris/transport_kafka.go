package gosiris

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

var Kafka = "kafka"

func init() {
	registerTransport(Kafka, newKafkaTransport)
}

func newKafkaTransport() TransportInterface {
	return new(kafkaTransport)
}

type kafkaTransport struct {
	url      string
	producer sarama.AsyncProducer
	consumer sarama.Consumer
}

type accessLogEntry struct {
	Method       string  `json:"method"`
	Host         string  `json:"host"`
	Path         string  `json:"path"`
	IP           string  `json:"ip"`
	ResponseTime float64 `json:"response_time"`

	encoded []byte
	err     error
}

func (k *kafkaTransport) Configure(url string, options map[string]string) {
	k.url = url
}

func (k *kafkaTransport) Connection() error {
	list := strings.Split(k.url, ",")

	producer, err := newProducer(list)
	if err != nil {
		return err
	}

	consumer, err := newConsumer(list)
	if err != nil {
		return err
	}

	k.producer = producer
	k.consumer = consumer

	return nil
}

func (k *kafkaTransport) Receive(queueName string) {
	consumer, err := k.consumer.ConsumePartition(queueName, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-consumer.Errors():
			ErrorLogger.Printf("Kafka consumer error: %v", err)
		case message := <-consumer.Messages():
			msg := EmptyContext
			json.Unmarshal(message.Value, &msg)
			InfoLogger.Printf("New Kafka message received: %v", msg)
			ActorSystem().Invoke(msg)
		}
	}
}

func (k *kafkaTransport) Close() {
	if err := k.producer.Close(); err != nil {
		ErrorLogger.Printf("Failed shut access log producer: %v", err)
	}
}

func (k *kafkaTransport) Send(destination string, data []byte) error {
	InfoLogger.Printf("Sending message to the Kafka destination %v", destination)
	k.producer.Input() <- &sarama.ProducerMessage{
		Topic: destination,
		Value: sarama.StringEncoder(data),
	}

	return nil
}

func newConsumer(brokerList []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 15 * time.Millisecond

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		ErrorLogger.Printf("Failed to start sarama consumer: %v", err)
		return nil, fmt.Errorf("failed to start sarama consumer: %v", err)
	}

	return consumer, nil
}

func newProducer(brokerList []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal      // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy  // Compress messages
	config.Producer.Flush.Frequency = 15 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		ErrorLogger.Printf("Failed to start sarama producer: %v", err)
		return nil, fmt.Errorf("failed to start sarama producer: %v", err)
	}

	go func() {
		for err := range producer.Errors() {
			ErrorLogger.Printf("Failed to write access log entry: %v", err)
		}
	}()

	return producer, nil
}
