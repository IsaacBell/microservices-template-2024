package server

import (
	helloworld "microservices-template-2024/api/helloworld/v1"
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

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

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UsersService, logger log.Logger) *http.Server {
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
	helloworld.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUsersHTTPServer(srv, user)

	StartPrometheus(srv)
	return srv
}
