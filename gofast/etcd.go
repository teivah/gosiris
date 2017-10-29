package gofast

import (
	"github.com/coreos/etcd/client"
	"time"
	"Gofast/gofast/util"
	"context"
	"strings"
)

const (
	actors_configuration = "/proactor/actor/"
	prefix               = "gofast://"
	delimiter            = "#"
)

type etcdClient struct {
	api client.KeysAPI
}

func (etcdClient *etcdClient) Configure(url ...string) error {
	cfg := client.Config{
		Endpoints:               url,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		util.LogError("etcd connection error %v", err)
		return err
	}
	etcdClient.api = client.NewKeysAPI(c)

	return nil
}

func (etcdClient *etcdClient) ParseConfiguration() (map[string]actorRemoteConfiguration, error) {
	resp, err := etcdClient.Get(actors_configuration)

	if err != nil {
		return nil, nil
	}

	conf := make(map[string]actorRemoteConfiguration)

	nodes := resp.Node.Nodes
	for i := 0; i < nodes.Len(); i++ {
		v := nodes[i].Value
		a :=
			strings.Split(v, delimiter)
		conf[nodes[i].Key] = actorRemoteConfiguration{a[0], a[1], a[2]}
	}

	return conf, nil
}

func (etcdClient *etcdClient) Set(key string, value string) error {
	_, err := etcdClient.api.Set(context.Background(), key, value, nil)

	if err != nil {
		util.LogError("etcd set %v error %v", key, err)
	}

	return err
}

func (etcdClient *etcdClient) Get(key string) (*client.Response, error) {
	resp, err := etcdClient.api.Get(context.Background(), key, &client.GetOptions{false, false, false})

	if err != nil {
		util.LogError("etcd get %v error %v", key, err)
		return resp, err
	}

	return resp, err
}

func (etcdClient *etcdClient) GetValue(key string) (string, error) {
	resp, err := etcdClient.Get(key)

	return resp.Node.Value, err
}
