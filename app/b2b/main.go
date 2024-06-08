package main

import (
	"core/internal/server"

	discovery_etcd "core/internal/discovery/etcd"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name     string = "leads" // name of the compiled software.
	Version  string           // version of the compiled software.
	flagconf string           // flagconf is the config flag.

	Watcher *discovery_etcd.Watcher // service discovery

	id string = server.GenerateServiceInstanceID()

	KafkaTopics = []string{"b2b", "leads", "leads/cdc", "companies"}
)

func init() {
	server.InitEnv(id, Name, &flagconf, KafkaTopics)
}

func newB2bApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	watcher, app := server.NewApp(Name, id, Version, logger, gs, hs)
	Watcher = watcher
	return app
}

func main() {
	server.RunApp(Name, Version, flagconf, Watcher, wireApp, nil)
}
