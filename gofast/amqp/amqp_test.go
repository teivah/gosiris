package amqp

import "testing"

func TestConnection(t *testing.T) {
	InitConfiguration("amqp://guest:guest@amqp:5672/")
}
