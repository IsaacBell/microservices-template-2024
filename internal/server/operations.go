package server

import (
	"context"
	"core/internal/conf"
	discovery_etcd "core/internal/discovery/etcd"
	"core/internal/util"
	"core/pkg/influx"
	"core/pkg/stream"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	etcdClient "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/etcd/registry"
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

var (
	flagconf string // config flag
	Watcher  *discovery_etcd.Watcher
	reg      *registry.Registry
)

func InitKafkaConsumer(serviceName string, topics []string) {
	for _, topic := range topics {
		stream.StartKafkaConsumer(topic, serviceName, func(msg string) {
			log.Infof("Kafka: [", topic, "] ", msg)
		})
		fmt.Println("consuming kafka topic: ", topic)
	}
}

func StartInfluxDb() {
	influx.TestInfluxV3()
}

func InitEnv(serviceName string, flagconf *string, topics []string) {
	file := conf.ConfigDir() + "config.yaml"
	flag.StringVar(flagconf, "conf", file, "config path, eg: -conf config.yaml")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory: ", err)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Current working directory: ", cwd)
		fmt.Println("err loading .env: ", err)

		err = godotenv.Load("../../.env")
		if err != nil {
			log.Fatalf("err loading .env: %v", err)
		}
	}

	InitKafkaConsumer(serviceName, topics)
	StartInfluxDb()

	// service discovery-level context
	ctx := AppCtx(serviceName)

	// Create etcd registrar
	etcdClient, regTmp := discovery_etcd.Register(ctx, serviceName)
	reg = regTmp
	Watcher, _ = discovery_etcd.NewWatcher(ctx, serviceName, etcdClient)
}

func watchServiceDiscovery(watcher *discovery_etcd.Watcher) {
	util.PrintLnInColor(util.AnsiColorBlue, "Running service discovery...")
	for {
		instances, err := watcher.Next()
		if err != nil {
			util.PrintLnInColor(util.AnsiColorRed, "Failed to get next instances:", err)
			if err == context.Canceled {
				// Watcher has been stopped, exit the loop
				break
			}
			time.Sleep(time.Second)
			continue
		}

		if len(instances) == 0 {
			util.PrintLnInColor(util.AnsiColorYellow, "No services discovered")
		} else {
			util.PrintLnInColor(util.AnsiColorMagenta, "::::::::Currently registered services::::::::")
			for _, instance := range instances {
				util.PrintLnInColor(util.AnsiColorCyan, "Service instance", util.AnsiColorGreen, instance)
			}
		}
		time.Sleep(time.Second)
	}
	util.PrintLnInColor(util.AnsiColorBlue, "Service discovery stopped.")
}

func StartServiceDiscovery(w *discovery_etcd.Watcher) {
	var wg sync.WaitGroup
	defer wg.Done()

	if w == nil {
		log.Fatalf("can't start service discovery: watcher not found")
		// return nil
	} else {
		wg.Add(1)
		go func() {
			w.StartDiscovery()
		}()
		wg.Wait()
	}
}

func AppCtx(serviceName string) context.Context {
	ctx := context.WithValue(context.Background(), "_app_service_name", serviceName)
	return context.WithValue(ctx, "_app_config", flagconf)
}

/*
Initializes an instance of the app.

Returns the app as well as an etcd watcher for service discovery.
*/
func NewApp(name, id, ver string, logger log.Logger, gs *grpc.Server, hs *http.Server) (*discovery_etcd.Watcher, *kratos.App) {
	fmt.Println("service name:", name)
	fmt.Println("machine user id:", id)

	return Watcher, kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(ver),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(reg),
	)
}

/*
 */
func RunApp(
	name, version, flagconf string,
	watcher *discovery_etcd.Watcher,
	wireAppFunc func(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error),
	afterStartCb func(),
) {
	var servers *conf.Server

	id, _ := os.Hostname()

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

	// Get config from the `configs/config.yaml` file
	// Config values must be set for each available service
	// Each service must have an http address and grpc address configured

	var configPath string
	if name == "core" {
		configPath = "server"
	} else {
		configPath = name
	}

	httpAddr, err := c.Value(configPath + ".http.addr").String()
	if err != nil {
		panic(err)
	}
	httpTimeout, err := c.Value(configPath + ".http.timeout").String()
	if err != nil {
		panic(err)
	}
	grpcAddr, err := c.Value(configPath + ".grpc.addr").String()
	if err != nil {
		panic(err)
	}
	grpcTimeout, err := c.Value(configPath + ".grpc.timeout").String()
	if err != nil {
		panic(err)
	}

	servers = &conf.Server{
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

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	db := &conf.Data_Database{
		Driver: "postgresql",
		Source: DbConnString(),
	}

	data := &conf.Data{
		Database: db,
		Redis:    bc.Data.Redis,
	}

	// Set redis to our url
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

	util.PrintLnInColor(util.AnsiColorGreen, "\n:::::", util.AnsiColorYellow, "service booting...", util.AnsiColorGreen, ":::::\n\n", name)
	defer util.PrintLnInColor(util.AnsiColorGreen, "\n:::::", util.AnsiColorYellow, "service shutting down...", util.AnsiColorGreen, ":::::\n\n", name)

	go util.RecordSystemMetrics()
	StartServiceDiscovery(Watcher)
	discovery_etcd.Watch(context.Background(), Watcher.Client, "/"+name, func(event *etcdClient.Event) {
		fmt.Println("------------------------")
		fmt.Printf("event.Kv: %v\n", event.Kv)
		fmt.Printf("event.Type: %v\n", event.Type)
		fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
		fmt.Println("------------------------")
	})
	discovery_etcd.Watch(context.Background(), Watcher.Client, "/example", func(event *etcdClient.Event) {
		fmt.Println("------------------------")
		fmt.Printf("event.Kv: %v\n", event.Kv)
		fmt.Printf("event.Type: %v\n", event.Type)
		fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
		fmt.Println("------------------------")
	})
	go discovery_etcd.Watch(context.Background(), Watcher.Client, "/example", func(event *etcdClient.Event) {
		fmt.Println("------------------------")
		fmt.Printf("event.Kv: %v\n", event.Kv)
		fmt.Printf("event.Type: %v\n", event.Type)
		fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
		fmt.Println("------------------------")
	})
	discovery_etcd.Watch(context.Background(), Watcher.Client, "/system", func(event *etcdClient.Event) {
		fmt.Println("------------------------")
		fmt.Printf("event.Kv: %v\n", event.Kv)
		fmt.Printf("event.Type: %v\n", event.Type)
		fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
		fmt.Println("------------------------")
	})
	go discovery_etcd.Watch(context.Background(), Watcher.Client, "/system", func(event *etcdClient.Event) {
		fmt.Println("------------------------")
		fmt.Printf("event.Kv: %v\n", event.Kv)
		fmt.Printf("event.Type: %v\n", event.Type)
		fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
		fmt.Println("------------------------")
	})

	util.PrintLnInColor(util.AnsiColorGray, "Starting service discovery...")

	kratos.AfterStart(func(context.Context) error {
		if watcher != nil {
			util.PrintLnInColor(util.AnsiColorGreen, "yes3!!")
			Watcher = watcher
			discovery_etcd.Watch(context.Background(), Watcher.Client, "/"+name, func(event *etcdClient.Event) {
				fmt.Println("------------------------")
				fmt.Printf("event.Kv: %v\n", event.Kv)
				fmt.Printf("event.Type: %v\n", event.Type)
				fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
				fmt.Println("------------------------")
			})
			go watcher.StartDiscovery()
		} else if Watcher != nil {
			util.PrintLnInColor(util.AnsiColorGreen, "yes2!!")
			discovery_etcd.Watch(context.Background(), Watcher.Client, "/"+name, func(event *etcdClient.Event) {
				fmt.Println("------------------------")
				fmt.Printf("event.Kv: %v\n", event.Kv)
				fmt.Printf("event.Type: %v\n", event.Type)
				fmt.Printf("event.PrevKv: %v\n", event.PrevKv)
				fmt.Println("------------------------")
			})
		} else {
			util.PrintLnInColor(util.AnsiColorYellow, "No watcher found after starting kratos")
		}
		stream.ProduceKafkaMessage(name, name+" server started")

		if afterStartCb != nil {
			afterStartCb()
		}

		return nil
	})
	kratos.AfterStop(func(context.Context) error {
		if watcher != nil {
			Watcher = watcher
			util.PrintLnInColor(util.AnsiColorMagenta, "Stopping service discovery watcher")
			defer watcher.Stop()
		}
		stream.ProduceKafkaMessage(name, name+" server de-registered")
		stream.ProduceKafkaMessage("system", name+" server de-registered")
		return nil
	})
	fmt.Printf("app.Run watcher: %v\n", Watcher)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

	defer app.Stop()
}
