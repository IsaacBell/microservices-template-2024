package leads_server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	prometheusClient "github.com/prometheus/client_golang/prometheus"
	promHttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Declare server types to run concurrently at runtime
var ProviderSet = wire.NewSet(
	NewLeadGrpcServer, NewLeadHTTPServer,
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
