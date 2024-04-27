package main

import (
	"flag"
	"fmt"
	"os"

	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	var file string
	_, err := os.Stat("./configs/config.yaml")
	if os.IsNotExist(err) {
		file = "../../configs/config.yaml"
	} else {
		file = "./configs/config.yaml"
	}
	flag.StringVar(&flagconf, "conf", file, "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Current working directory: %s", cwd)
		fmt.Println("err loading .env: %v", err)

		err = godotenv.Load("../../.env")
		if err != nil {
			log.Fatalf("err loading .env: %v", err)
		}
	}

	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	bc.Data.Database = &conf.Data_Database{
		Driver: "postgresql",
		Source: server.DbConnString(),
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := server.OpenDBConn(); err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
