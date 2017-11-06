package gosiris

import (
	"fmt"
)

var remoteConnections map[string]RemoteConnectionInterface

type RemoteConnectionInterface interface {
	Configure(string)
	Connection() error
	Send(string, []byte)
	Receive(string)
	Close()
}

func InitRemoteConnections(configuration map[string]OptionsInterface) {
	remoteConnections = make(map[string]RemoteConnectionInterface)

	//TODO Dynamic management
	for k, v := range configuration {
		if v.RemoteType() == "amqp" {
			c := AmqpConnection{}
			c.Configure(v.Url())
			err := c.Connection()
			if err != nil {
				ErrorLogger.Printf("Failed to initialize the connection with %v: %v", v, err)
			}
			remoteConnections[k] = &c
		}
	}

	InfoLogger.Printf("Remote connections: %v", remoteConnections)
}

func AddConnection(name string, conf OptionsInterface) {
	if conf.RemoteType() == "amqp" {
		c := AmqpConnection{}
		c.Configure(conf.Url())
		err := c.Connection()
		if err != nil {
			ErrorLogger.Printf("Failed to initialize the connection with %v: %v", name, err)
		}
		remoteConnections[name] = &c
		InfoLogger.Printf("AMQP %v connection added", name)
	}
}

func DeleteRemoteActorConnection(name string) error {
	v, exists := remoteConnections[name]
	if !exists {
		ErrorLogger.Printf("Delete error: connection %v not registered", name)
		return fmt.Errorf("delete error: connection %v not registered", name)
	}

	v.Close()
	delete(remoteConnections, name)

	InfoLogger.Printf("Connection %v deleted", name)

	return nil
}

func RemoteConnection(name string) (RemoteConnectionInterface, error) {
	v, exists := remoteConnections[name]
	if !exists {
		ErrorLogger.Printf("Remote connection error: connection %v not registered", name)
		return nil, fmt.Errorf("remote connection error: connection %v not registered", name)
	}

	return v, nil
}
