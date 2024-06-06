//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"core/internal/conf"
	leads_biz "core/pkg/leads/biz"
	leads_data "core/pkg/leads/data"
	leads_server "core/pkg/leads/server"
	leads_service "core/pkg/leads/service"

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
