package leads_server

import (
	leadsV1 "microservices-template-2024/api/v1/b2b"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/server"
	leadService "microservices-template-2024/pkg/leads/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewLeadGrpcServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	lead *leadService.LeadService,
	company *leadService.CompanyService,
) *grpc.Server {
	srv := server.GRPCServerFactory("b2b", c, logger)
	leadsV1.RegisterLeadsServer(srv, lead)
	leadsV1.RegisterCompaniesServer(srv, company)

	return srv
}

func NewLeadHTTPServer(
	c *conf.Server,
	logger log.Logger,
	// Each available runtime service
	lead *leadService.LeadService,
	company *leadService.CompanyService,
) *http.Server {
	srv := server.HTTPServerFactory("b2b", c, logger)
	leadsV1.RegisterLeadsHTTPServer(srv, lead)
	leadsV1.RegisterCompaniesHTTPServer(srv, company)

	server.StartPrometheus(srv)
	return srv
}
