package gosiris

import (
	"testing"
)

func TestAmqpUnit(t *testing.T) {
	r := amqpTransport{}
	r.Configure("amqp://guest:guest@amqp:5672/", nil)

	r.Send("test", nil)

	r.Receive("test")
}
