package gosiris

import (
	"testing"
)

func TestPublish(t *testing.T) {
	r := AmqpConnection{}
	r.Configure("amqp://guest:guest@amqp:5672/", nil)

	r.Send("test", nil)

	r.Receive("test")
}
