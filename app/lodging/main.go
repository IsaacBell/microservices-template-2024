package main

import (
	"flag"
	"fmt"
	"os"

	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/server"
	stream "microservices-template-2024/pkg/stream"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/types/known/durationpb"
	// "google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "Lodging API"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()

	KafkaTopics = []string{"lodging", "properties", "properties/cdc"}
)

func init() {
	file := conf.ConfigDir() + "config.yaml"
	flag.StringVar(&flagconf, "conf", file, "config path, eg: -conf config.yaml")

	streamKafkaMessages()
}

func newLodgingApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	fmt.Println("service name:", Name)
	fmt.Println("machine user id:", id)
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

func streamKafkaMessages() {
	for _, topic := range KafkaTopics {
		stream.StartKafkaConsumer(topic, "core", func(msg string) {
			log.Infof("Kafka: [", topic, "] ", msg)
		})
	}
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
	fmt.Println("flag", flagconf)
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

	http_addr, err := c.Value("leads.http.addr").String()
	if err != nil {
		panic(err)
	}
	http_timeout, err := c.Value("leads.http.timeout").String()
	if err != nil {
		panic(err)
	}
	grpc_addr, err := c.Value("leads.grpc.addr").String()
	if err != nil {
		panic(err)
	}
	grpc_timeout, err := c.Value("leads.grpc.timeout").String()
	if err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	servers := &conf.Server{
		Grpc: &conf.Server_GRPC{
			Network: "",
			Addr:    grpc_addr,
			Timeout: &durationpb.Duration{Seconds: int64(grpc_timeout[0])},
		},
		Http: &conf.Server_HTTP{
			Network: "",
			Addr:    http_addr,
			Timeout: &durationpb.Duration{Seconds: int64(http_timeout[0])},
		},
	}

	db := &conf.Data_Database{
		Driver: "postgresql",
		Source: server.DbConnString(),
	}

	data := &conf.Data{
		Database: db,
		Redis:    bc.Data.Redis,
	}

	pass := os.Getenv("UPSTASH_REDIS_PASS")
	url := os.Getenv("UPSTASH_REDIS_URL")
	path := "rediss://default:" + pass + "@" + url
	bc.Data.Redis.Addr = path
	fmt.Println(servers.String(), data.String(), logger.Log)

	app, cleanup, err := wireApp(servers, data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := server.OpenDBConn(); err != nil {
		panic(err)
	}

	fmt.Println("::::: Lodging API online :::::")
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
	defer app.Stop()

}
