package gosiris

import (
	"fmt"
)

var remoteConnections map[string]TransportInterface
var transportTypes map[string]func() TransportInterface

func init() {
	transportTypes = make(map[string]func() TransportInterface)
}

type TransportInterface interface {
	Configure(string, map[string]string)
	Connection() error
	Send(string, []byte) error
	Receive(string)
	Close()
}

func registerTransport(name string, f func() TransportInterface) {
	transportTypes[name] = f
}

func newTransport(name string) TransportInterface {
	return transportTypes[name]()
}

func InitRemoteConnections(configuration map[string]OptionsInterface) {
	remoteConnections = make(map[string]TransportInterface)

	for k, v := range configuration {
		c := transportTypes[v.RemoteType()]()

		c.Configure(v.Url(), nil)
		err := c.Connection()
		if err != nil {
			ErrorLogger.Printf("Failed to initialize the connection with %v: %v", v, err)
		}
		remoteConnections[k] = c
	}

	InfoLogger.Printf("Remote connections: %v", remoteConnections)
}

func AddConnection(name string, conf OptionsInterface) {
	c := transportTypes[conf.RemoteType()]()

	c.Configure(conf.Url(), nil)
	err := c.Connection()
	if err != nil {
		ErrorLogger.Printf("Failed to initialize the connection with %v: %v", name, err)
	}
	remoteConnections[name] = c
	InfoLogger.Printf("Remote connection %v added", name)
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

func RemoteConnection(name string) (TransportInterface, error) {
	v, exists := remoteConnections[name]
	if !exists {
		ErrorLogger.Printf("Remote connection error: connection %v not registered", name)
		return nil, fmt.Errorf("remote connection error: connection %v not registered", name)
	}

	return v, nil
}
