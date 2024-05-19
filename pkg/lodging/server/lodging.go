package lodging_server

import (
	"fmt"
	"io/ioutil"
	lodgingV1 "microservices-template-2024/api/v1/lodging"
	"microservices-template-2024/internal/conf"
	lodgingService "microservices-template-2024/pkg/lodging/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/go-kratos/kratos/v2/middleware/metrics"
	kmetrics "github.com/go-kratos/prometheus/metrics"
)

func NewLodgingGrpcServer(
	c *conf.Server,
	logger log.Logger,
	property *lodgingService.PropertyService,
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
	lodgingV1.RegisterLodgingServer(srv, property)

	return srv
}

func NewLodgingHTTPServer(
	c *conf.Server,
	logger log.Logger,
	property *lodgingService.PropertyService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	lodgingV1.RegisterLodgingHTTPServer(srv, property)

	StartPrometheus(srv)
	return srv
}
