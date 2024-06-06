package finance_server

import (
	v1 "core/api/v1"
	"core/internal/conf"
	"core/internal/server"
	finance_service "core/pkg/finance/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewFinanceGrpcServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	fin *finance_service.FinanceService,
) *grpc.Server {
	srv := server.GRPCServerFactory("finance", c, logger)
	v1.RegisterFinanceServer(srv, fin)

	return srv
}

func NewFinanceHTTPServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	fin *finance_service.FinanceService,
) *http.Server {
	srv := server.HTTPServerFactory("finance", c, logger)
	v1.RegisterFinanceHTTPServer(srv, fin)

	server.StartPrometheus(srv)
	return srv
}
