package gopera

import (
	"time"
	"math/rand"
	"testing"
	"gopera/gopera/util"
)

var c etcdClient

func init() {
	c = etcdClient{}
	c.Configure("http://etcd:2379")
}

func randomString(size int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestRandom(t *testing.T) {
	k := randomString(12)
	v := randomString(12)
	c.Set(k, v)

	value, err := c.GetValue(k)

	if err != nil {
		t.Errorf("Get error %v", err)
		t.Fail()
		return
	}

	if value != v {
		t.Errorf("The value is not equals")
		t.Fail()
		return
	}
}

func TestFixed(t *testing.T) {
	k := "/proactor/actor/actor2"
	v := "amqp#amqp://guest:guest@amqp:5672/#actor2"
	c.Set(k, v)

	value, err := c.GetValue(k)

	if err != nil {
		t.Errorf("Get error %v", err)
		t.Fail()
		return
	}

	if value != v {
		t.Errorf("The value is not equals")
		t.Fail()
		return
	}
}

func TestGet(t *testing.T) {
	value, _ := c.ParseConfiguration()

	util.LogInfo("conf: %v", value)
}

func TestWatch(t *testing.T) {
	c.Watch()
}