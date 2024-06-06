//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"core/internal/conf"
	consultants_biz "core/pkg/consultants/biz"
	consultants_data "core/pkg/consultants/data"
	consultants_server "core/pkg/consultants/server"
	consultants_service "core/pkg/consultants/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		consultants_server.ProviderSet, consultants_data.ProviderSet,
		consultants_biz.ProviderSet, consultants_service.ProviderSet, newConsultantsApp,
	))
}
