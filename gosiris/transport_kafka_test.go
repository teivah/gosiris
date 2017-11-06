package gosiris

import (
	"testing"
)

func TestKafkaUnit(t *testing.T) {
	r := kafkaTransport{}
	r.Configure("kafka:9092", nil)
	err := r.Connection()
	if err != nil {
		t.Errorf("Connection error: %v", err)
		t.Fail()
		return
	}
}
