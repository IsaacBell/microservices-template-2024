//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"microservices-template-2024/internal/conf"
	finance_biz "microservices-template-2024/pkg/finance/biz"
	finance_data "microservices-template-2024/pkg/finance/data"
	finance_server "microservices-template-2024/pkg/finance/server"
	finance_service "microservices-template-2024/pkg/finance/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		finance_server.ProviderSet, finance_data.ProviderSet,
		finance_biz.ProviderSet, finance_service.ProviderSet, newFinanceApp,
	))
}
