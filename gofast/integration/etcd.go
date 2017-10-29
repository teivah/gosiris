package integration

import (
	"github.com/coreos/etcd/client"
	"fmt"
	"time"
)

func InitConfiguration(endpoints ...string) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	fmt.Print(cfg)
}
