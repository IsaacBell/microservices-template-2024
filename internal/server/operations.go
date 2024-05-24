package server

import (
	"context"
	"flag"
	"fmt"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/pkg/stream"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/types/known/durationpb"
)

func InitKafkaConsumer(serviceName string, topics []string) {
	for _, topic := range topics {
		stream.StartKafkaConsumer(topic, serviceName, func(msg string) {
			log.Infof("Kafka: [", topic, "] ", msg)
		})
		fmt.Println("consuming kafka topic: ", topic)
	}
}

func InitEnv(serviceName, flagconf string, topics []string) {
	file := conf.ConfigDir() + "config.yaml"
	flag.StringVar(&flagconf, "conf", file, "config path, eg: -conf config.yaml")

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

	InitKafkaConsumer(serviceName, topics)
}

func NewApp(name, id, ver string, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	fmt.Println("service name:", name)
	fmt.Println("machine user id:", id)
	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(ver),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func RunApp(
	name, version, flagconf string,
	wireAppFunc func(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error),
	afterStartCb func(),
) {
	id, _ := os.Hostname()

	flag.Parse()
	fmt.Printf("flag: %v\n", flagconf)
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version,
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

	httpAddr, err := c.Value(name + ".http.addr").String()
	if err != nil {
		panic(err)
	}
	httpTimeout, err := c.Value(name + ".http.timeout").String()
	if err != nil {
		panic(err)
	}
	grpcAddr, err := c.Value(name + ".grpc.addr").String()
	if err != nil {
		panic(err)
	}
	grpcTimeout, err := c.Value(name + ".grpc.timeout").String()
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
			Addr:    grpcAddr,
			Timeout: &durationpb.Duration{Seconds: int64(grpcTimeout[0])},
		},
		Http: &conf.Server_HTTP{
			Network: "",
			Addr:    httpAddr,
			Timeout: &durationpb.Duration{Seconds: int64(httpTimeout[0])},
		},
	}

	db := &conf.Data_Database{
		Driver: "postgresql",
		Source: DbConnString(),
	}

	data := &conf.Data{
		Database: db,
		Redis:    bc.Data.Redis,
	}

	pass := os.Getenv("UPSTASH_REDIS_PASS")
	url := os.Getenv("UPSTASH_REDIS_URL")
	path := "rediss://default:" + pass + "@" + url
	bc.Data.Redis.Addr = path
	fmt.Printf("%s %s %v\n", servers.String(), data.String(), logger.Log)

	app, cleanup, err := wireAppFunc(servers, data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := OpenDBConn(); err != nil {
		panic(err)
	}

	fmt.Printf("::::: %s Service booting :::::\n", name)
	defer fmt.Printf("::::: %s Service shutting down :::::\n", name)

	kratos.AfterStart(func(context.Context) error {
		stream.ProduceKafkaMessage("main", name+"Finance Server started")
		defer stream.ProduceKafkaMessage("main", name+"Finance Server stopped")
		stream.ProduceKafkaMessage(name, name+" server started")
		defer stream.ProduceKafkaMessage(name, name+" server stopped")

		if afterStartCb != nil {
			afterStartCb()
		}

		return nil
	})

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
	defer app.Stop()
}
