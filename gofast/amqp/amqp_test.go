package amqp

import (
	"testing"
)

func TestPublish(t *testing.T) {
	r := DistributedAmqp{}
	r.Init("amqp", "amqp://guest:guest@amqp:5672/")

	r.Send("test")

	r.Receive("test")
}
