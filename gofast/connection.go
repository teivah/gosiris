package gofast

import (
	"fmt"
	"Gofast/gofast/util"
)

var manager map[string]RemoteConnectionInterface

type RemoteConnectionInterface interface {
	Configure(string)
	Connection() error
	Send(string, []byte)
	Receive(string)
	Close()
}

func InitRemoteConnections(configuration map[string]OptionsInterface) {
	fmt.Println("b")
	manager = make(map[string]RemoteConnectionInterface)

	//TODO Dynamic management
	for k, v := range configuration {
		if v.RemoteType() == "amqp" {
			c := AmqpConnection{}
			c.Configure(v.Url())
			err := c.Connection()
			if err != nil {
				util.LogError("Failed to initialize the connection with %v: %v", v, err)
			}
			manager[k] = &c
		}
	}
}

func AddRemoteConnection(name string, conf OptionsInterface) {
	if conf.RemoteType() == "amqp" {
		c := AmqpConnection{}
		c.Configure(conf.Url())
		err := c.Connection()
		if err != nil {
			util.LogError("Failed to initialize the connection with %v: %v", name, err)
		}
		manager[name] = &c
		util.LogInfo("AMQP %v connection added", name)
	}
}

func DeleteConnection(name string) error {
	v, exists := manager[name]
	if !exists {
		util.LogError("Connection %v not registered", name)
		return fmt.Errorf("connection %v not registered", name)
	}

	v.Close()
	delete(manager, name)

	util.LogInfo("Connection %v deleted")

	return nil
}
