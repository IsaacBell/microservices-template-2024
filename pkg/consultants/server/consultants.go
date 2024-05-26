package consultants_server

import (
	consultantsV1 "microservices-template-2024/api/v1/consultants"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/server"
	consultantsService "microservices-template-2024/pkg/consultants/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewConsultantsGrpcServer(
	c *conf.Server,
	logger log.Logger,
	consultant *consultantsService.ConsultantService,
) *grpc.Server {
	srv := server.GRPCServerFactory("consultants", c, logger)
	consultantsV1.RegisterConsultantsServer(srv, consultant)

	return srv
}

func NewConsultantsHTTPServer(
	c *conf.Server,
	logger log.Logger,
	consultant *consultantsService.ConsultantService,
) *http.Server {
	srv := server.HTTPServerFactory("consultants", c, logger)
	consultantsV1.RegisterConsultantsHTTPServer(srv, consultant)

	server.StartPrometheus(srv)
	return srv
}
