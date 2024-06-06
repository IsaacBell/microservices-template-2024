package server

import (
	helloworld "core/api/helloworld/v1"
	v1 "core/api/v1"
	"core/internal/auth"
	"core/internal/conf"
	"core/internal/service"
	analyticsengine "core/pkg/analyticsEngine"

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

func HTTPServerFactory(name string, c *conf.Server, logger log.Logger) *http.Server {
	authCtx := auth.NewAuthCtx()
	
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			auth.JwtMiddleware(authCtx), // must come before analytics
			analyticsengine.MoesifMiddleware(authCtx),
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

	return http.NewServer(opts...)
}

func NewCoreHTTPServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	greeter *service.GreeterService,
	user *service.UsersService,
	trans *service.TransactionsService,
	lias *service.LiabilitiesService,
	log *service.LogService,
) *http.Server {
	srv := HTTPServerFactory("core", c, logger)
	helloworld.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUsersHTTPServer(srv, user)
	v1.RegisterTransactionsHTTPServer(srv, trans)
	v1.RegisterLiabilitiesHTTPServer(srv, lias)
	v1.RegisterLogHTTPServer(srv, log)

	StartPrometheus(srv)
	return srv
}
