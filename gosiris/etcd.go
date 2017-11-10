package gosiris

import (
	"context"
	"github.com/coreos/etcd/client"
	"strings"
	"time"
)

const (
	actors_configuration = "/gosiris/actor/"
	prefix               = "gosiris://"
	delimiter            = "#"
	action_delete        = "delete"
	action_set           = "set"
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
		ErrorLogger.Printf("etcd connection error %v", err)
		return err
	}
	etcdClient.api = client.NewKeysAPI(c)
	etcdClient.createDir(actors_configuration)

	return nil
}

func (etcdClient *etcdClient) Close() {
}

func (etcdClient *etcdClient) Watch(cbCreate func(string, *ActorOptions), cbDelete func(string)) error {
	w := etcdClient.api.Watcher(actors_configuration, &client.WatcherOptions{
		AfterIndex: 0,
		Recursive:  true,
	})

	for {
		r, err := w.Next(context.Background())

		if err != nil {
			ErrorLogger.Printf("etcd watch error: %v", err)
			return err
		}

		if r.Action == action_delete {
			k := r.Node.Key
			InfoLogger.Printf("Actor %v removed from the registry", k)

			cbDelete(k)
		} else if r.Action == action_set {
			k, v := parseNode(r.Node)

			InfoLogger.Printf("New actor %v added to the registry", k)

			cbCreate(k, v)
		}
	}
}

func parseNode(node *client.Node) (string, *ActorOptions) {
	v := node.Value
	a :=
		strings.Split(v, delimiter)
	k := node.Key

	return k[len(actors_configuration):], &ActorOptions{a[0], true, true, a[1], a[2], a[3], 0, 0}
}

func (etcdClient *etcdClient) ParseConfiguration() (map[string]OptionsInterface, error) {
	resp, err := etcdClient.Get(actors_configuration)

	if err != nil {
		return nil, nil
	}

	conf := make(map[string]OptionsInterface)

	nodes := resp.Node.Nodes
	for i := 0; i < nodes.Len(); i++ {
		k, v := parseNode(nodes[i])
		conf[k] = v
	}

	return conf, nil
}

func (etcdClient *etcdClient) RegisterActor(name string, options OptionsInterface) error {
	k := actors_configuration + name
	v := options.Parent() + delimiter + options.RemoteType() + delimiter + options.Url() + delimiter + options.Destination()

	err := etcdClient.Set(k, v)
	if err != nil {
		ErrorLogger.Printf("Failed to register actor %v: %v", k, err)
		return err
	}

	return nil
}

func (etcdClient *etcdClient) UnregisterActor(name string) error {
	return etcdClient.Delete(actors_configuration + name)
}

func (etcdClient *etcdClient) createDir(key string) {
	opt := new(client.SetOptions)
	opt.Dir = true

	etcdClient.api.Set(context.Background(), key, "", opt)
}

func (etcdClient *etcdClient) Set(key string, value string) error {
	_, err := etcdClient.api.Set(context.Background(), key, value, nil)

	if err != nil {
		ErrorLogger.Printf("etcd set %v error %v", key, err)
	}

	return err
}

func (etcdClient *etcdClient) Delete(key string) error {
	_, err := etcdClient.api.Delete(context.Background(), key, nil)

	return err
}

func (etcdClient *etcdClient) Get(key string) (*client.Response, error) {
	resp, err := etcdClient.api.Get(context.Background(), key, &client.GetOptions{
		Recursive: false,
		Sort:      false,
		Quorum:    false,
	})

	if err != nil {
		ErrorLogger.Printf("etcd get %v error %v", key, err)
		return resp, err
	}

	return resp, err
}

func (etcdClient *etcdClient) GetValue(key string) (string, error) {
	resp, err := etcdClient.Get(key)

	return resp.Node.Value, err
}
