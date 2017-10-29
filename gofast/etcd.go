package gofast

import (
	"github.com/coreos/etcd/client"
	"time"
	"Gofast/gofast/util"
	"context"
	"strings"
	"fmt"
)

const (
	actors_configuration = "/proactor/actor/"
	prefix               = "gofast://"
	delimiter            = "#"
)

type etcdClient struct {
	api           client.KeysAPI
	configuration map[string]OptionsInterface
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

	etcdClient.configuration = make(map[string]OptionsInterface)

	return nil
}

func (etcdClient *etcdClient) ParseConfiguration() (map[string]OptionsInterface, error) {
	resp, err := etcdClient.Get(actors_configuration)

	if err != nil {
		return nil, nil
	}

	conf := make(map[string]OptionsInterface)

	nodes := resp.Node.Nodes
	for i := 0; i < nodes.Len(); i++ {
		v := nodes[i].Value
		a :=
			strings.Split(v, delimiter)
		conf[nodes[i].Key] = &ActorOptions{true, a[0], a[1], a[2]}
	}

	return conf, nil
}

func (etcdClient *etcdClient) RegisterActor(name string, options OptionsInterface) error {
	k := actors_configuration + name
	v := options.RemoteType() + delimiter + options.Url() + delimiter + options.Destination()

	err := etcdClient.Set(k, v)
	if err != nil {
		fmt.Errorf("Failed to register actor %v: %v", k, err)
		return err
	}

	//TODO Implement a watcher
	etcdClient.configuration[k] = &ActorOptions{true, options.RemoteType(), options.Url(), options.Destination()}
	return nil
}

func (etcdClient *etcdClient) UnregisterActor(name string) error {
	return etcdClient.Delete(actors_configuration + name)
}

func (etcdClient *etcdClient) Configuration() (map[string]OptionsInterface) {
	return etcdClient.configuration
}

func (etcdClient *etcdClient) Set(key string, value string) error {
	_, err := etcdClient.api.Set(context.Background(), key, value, nil)

	if err != nil {
		util.LogError("etcd set %v error %v", key, err)
	}

	return err
}

func (etcdClient *etcdClient) Delete(key string) error {
	_, err := etcdClient.api.Delete(context.Background(), key, nil)
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
