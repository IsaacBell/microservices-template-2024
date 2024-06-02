package main

import (
	"os"

	"microservices-template-2024/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name string = "core" // name of the compiled software.
	Version string 
	flagconf string // config flag.

	id, _ = os.Hostname()

	KafkaTopics = []string{"core", "default", "critical", "main"}
)

func init() {
	server.InitEnv(Name, &flagconf, KafkaTopics)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return server.NewApp(Name, id, Version, logger, gs, hs)
}

func main() {
	server.RunApp(Name, Version, flagconf, wireApp, nil)
}
