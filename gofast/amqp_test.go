package gofast

import (
	"testing"
)

func TestPublish(t *testing.T) {
	r := DistributedAmqp{}
	r.Configure("amqp", "amqp://guest:guest@amqp:5672/")

	r.Send("test", nil)

	r.Receive("test")
}
