package lodging_server

import (
	lodgingV1 "core/api/v1/lodging"
	"core/internal/conf"
	"core/internal/server"
	lodgingService "core/pkg/lodging/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewLodgingGrpcServer(
	c *conf.Server,
	logger log.Logger,
	property *lodgingService.PropertyService,
) *grpc.Server {
	srv := server.GRPCServerFactory("lodging", c, logger)
	lodgingV1.RegisterLodgingServer(srv, property)

	return srv
}

func NewLodgingHTTPServer(
	c *conf.Server,
	logger log.Logger,
	property *lodgingService.PropertyService,
) *http.Server {
	srv := server.HTTPServerFactory("lodging", c, logger)
	lodgingV1.RegisterLodgingHTTPServer(srv, property)

	server.StartPrometheus(srv)
	return srv
}
