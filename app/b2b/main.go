package main

import (
	"os"

	"core/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	// "google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "leads"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()

	KafkaTopics = []string{"b2b", "leads", "leads/cdc", "companies"}
)

func init() {
	server.InitEnv(Name, &flagconf, KafkaTopics)
}

func newB2bApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return server.NewApp(Name, id, Version, logger, gs, hs)
}

func main() {
	server.RunApp(Name, Version, flagconf, wireApp, nil)
}
