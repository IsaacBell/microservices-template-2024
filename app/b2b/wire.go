//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"microservices-template-2024/internal/conf"
	leads_biz "microservices-template-2024/pkg/leads/biz"
	leads_data "microservices-template-2024/pkg/leads/data"
	leads_server "microservices-template-2024/pkg/leads/server"
	leads_service "microservices-template-2024/pkg/leads/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		leads_server.ProviderSet, leads_data.ProviderSet,
		leads_biz.ProviderSet, leads_service.ProviderSet, newB2bApp,
	))
}
