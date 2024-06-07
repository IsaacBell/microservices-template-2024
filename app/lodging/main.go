package main

import (
	"os"

	"core/internal/server"

	discovery_etcd "core/internal/discovery/etcd"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name     string                  = "lodging" // name of the compiled software.
	Version  string                              // version of the compiled software.
	flagconf string                              // flagconf is the config flag.
	Watcher  *discovery_etcd.Watcher             // service discovery

	id, _ = os.Hostname()

	KafkaTopics = []string{"lodging", "properties", "properties/cdc"}
)

func init() {
	server.InitEnv(Name, &flagconf, KafkaTopics)
}

func newLodgingApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	watcher, app := server.NewApp(Name, id, Version, logger, gs, hs)
	Watcher = watcher
	return app
}

func main() {
	server.RunApp(Name, Version, flagconf, Watcher, wireApp, nil)
}
