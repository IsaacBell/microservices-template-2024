package main

import (
	discovery_etcd "core/internal/discovery/etcd"
	zap "core/internal/logs"
	"core/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name     string = "core" // name of the compiled software.
	Version  string
	flagconf string // config flag.
	Log      zap.Logger
	Watcher  *discovery_etcd.Watcher // service discovery

	id string = server.GenerateServiceInstanceID()

	KafkaTopics = []string{"core", "default", "critical", "main"}
)

func init() {
	server.InitEnv(id, Name, &flagconf, KafkaTopics)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	w, app := server.NewApp(Name, id, Version, logger, gs, hs)
	Watcher = w
	return app
}

func main() {
	server.RunApp(Name, Version, flagconf, Watcher, wireApp, nil)
}
