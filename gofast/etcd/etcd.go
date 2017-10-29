package etcd

import (
	"github.com/coreos/etcd/client"
	"time"
	"Gofast/gofast/util"
	"context"
)

var api client.KeysAPI

func InitConfiguration(endpoints ...string) error {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		util.LogError("etcd connection error %v", err)
		return err
	}
	api = client.NewKeysAPI(c)

	return nil
}

func Set(key string, value string) error {
	_, err := api.Set(context.Background(), key, value, nil)

	if err != nil {
		util.LogError("etcd set %v error %v", key, err)
	}

	return err
}

func Get(key string) (string, error) {
	resp, err := api.Get(context.Background(), key, nil)

	if err != nil {
		util.LogError("etcd get %v error %v", key, err)
		return "", err
	}

	return resp.Node.Value, nil
}
