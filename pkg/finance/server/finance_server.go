package finance_server

import (
	"fmt"
	"io/ioutil"
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/conf"
	finance_service "microservices-template-2024/pkg/finance/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/go-kratos/kratos/v2/middleware/metrics"
	kmetrics "github.com/go-kratos/prometheus/metrics"

	"github.com/prometheus/client_golang/prometheus"
	prometheusClient "github.com/prometheus/client_golang/prometheus"
	promHttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheus(srv *http.Server) {
	_metricSeconds := prometheusClient.NewHistogramVec(prometheusClient.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duratio(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricRequests := prometheusClient.NewCounterVec(prometheusClient.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})

	prometheusClient.MustRegister(_metricSeconds, _metricRequests)

	srv.Handle("/metrics", promHttp.Handler())
}

func NewFinanceGrpcServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	fin *finance_service.FinanceService,
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
	v1.RegisterFinanceServer(srv, fin)

	return srv
}

func NewFinanceHTTPServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	fin *finance_service.FinanceService,
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
	v1.RegisterFinanceHTTPServer(srv, fin)

	StartPrometheus(srv)
	return srv
}
