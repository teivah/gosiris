package gosiris

import (
	"fmt"

	"github.com/couchbase/gocb"
)

var Couchbase = "couchbase"

type couchbaseTransport struct {
	url        string
	connection *gocb.Cluster
	bucket     *gocb.Bucket
}

func init() {
	registerTransport(Couchbase, newCouchbaseTransport)
}

func newCouchbaseTransport() TransportInterface {
	return new(couchbaseTransport)
}

func (a *couchbaseTransport) Configure(url string, options map[string]string) {
	a.url = url
}

func (a *couchbaseTransport) Connection() error {
	c, err := gocb.Connect(a.url)
	if err != nil {
		ErrorLogger.Printf("Failed to connect to the Couchbase server %v", a.url)
		return err
	}
	a.connection = c

	bucket, err := c.OpenBucket("travel-sample", "")
	if err != nil {
		ErrorLogger.Printf("Failed to open the travel-sample bucket on the server %v", a.url)
		return err
	}
	bucket.Manager("", "").CreatePrimaryIndex("", true, false)
	a.bucket = bucket

	InfoLogger.Printf("Connected to %v", a.url)

	return nil
}

type Data struct {
	data string `json:"data"`
}

func (a *couchbaseTransport) Receive(docId string) {
	var msg Context
	a.bucket.Get(docId, &msg)

	fmt.Printf("docId = %+v inData = %+v\n", docId, msg)

	InfoLogger.Printf("New couchbase message received: %v", msg)
	ActorSystem().Invoke(msg)
}

func (a *couchbaseTransport) Close() {
	a.bucket.Close()
}

func (a *couchbaseTransport) Send(destination string, data []byte) error {
	InfoLogger.Printf("Sending message %v to the couchbase destination %v", string(data), destination)

	body := data
	_, err := a.bucket.Upsert(
		destination, // exchange
		//Data{string(body)}, 0)
		string(body), 0)

	if err != nil {
		ErrorLogger.Printf("Error while inserting a message to couchbase bucket %v: %v", destination, err)
		return err
	}

	return nil
}
