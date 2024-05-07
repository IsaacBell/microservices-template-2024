package server

import (
	"fmt"
	"io/ioutil"
	helloworld "microservices-template-2024/api/helloworld/v1"
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/go-kratos/kratos/v2/middleware/metrics"
	kmetrics "github.com/go-kratos/prometheus/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewGRPCServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	greeter *service.GreeterService,
	user *service.UsersService,
	trans *service.TransactionsService,
	lias *service.LiabilitiesService,
	log *service.LogService,
) *grpc.Server {
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(ioutil.Discard))
	if err != nil {
		fmt.Println("creating stdout exporter: %v", err)
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("trace")),
		),
	)

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "kratos_counter"}, []string{"kind", "operation", "code", "reason"})

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(tracing.WithTracerProvider(tp)),
			metrics.Server(metrics.WithRequests(kmetrics.NewCounter(counter))),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(srv, greeter)
	v1.RegisterUsersServer(srv, user)
	v1.RegisterTransactionsServer(srv, trans)
	v1.RegisterLiabilitiesServer(srv, lias)
	v1.RegisterLogServer(srv, log)
	return srv
}
