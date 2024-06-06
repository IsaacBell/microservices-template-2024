//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"core/internal/conf"
	finance_biz "core/pkg/finance/biz"
	finance_data "core/pkg/finance/data"
	finance_server "core/pkg/finance/server"
	finance_service "core/pkg/finance/service"

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
