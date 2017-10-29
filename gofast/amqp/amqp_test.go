package amqp

import (
	"testing"
)

func TestPublish(t *testing.T) {
	InitConfiguration("amqp://guest:guest@amqp:5672/")

	Publish("test")

	AddConsumer("test")
}
