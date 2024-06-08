package main

import (
	discovery_etcd "core/internal/discovery/etcd"
	"core/internal/server"
	finance_util "core/pkg/finance/util"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name     string                  = "finance" // name of the compiled software.
	Version  string                              // version of the compiled software.
	flagconf string                              // flagconf is the config flag.
	Watcher  *discovery_etcd.Watcher             // service discovery

	id string = server.GenerateServiceInstanceID()

	KafkaTopics = []string{"finance", "finance/cdc"}
)

func init() {
	server.InitEnv(id, Name, &flagconf, KafkaTopics)
}

func newFinanceApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	watcher, app := server.NewApp(Name, id, Version, logger, gs, hs)
	Watcher = watcher
	return app
}

func main() {
	server.RunApp(Name, Version, flagconf, Watcher, wireApp, finance_util.InitFinnhubClient)
}
