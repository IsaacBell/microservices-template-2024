//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"core/internal/conf"
	lodging_biz "core/pkg/lodging/biz"
	lodging_data "core/pkg/lodging/data"
	lodging_server "core/pkg/lodging/server"
	lodging_service "core/pkg/lodging/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		lodging_server.ProviderSet, lodging_data.ProviderSet,
		lodging_biz.ProviderSet, lodging_service.ProviderSet, newLodgingApp,
	))
}
