//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"microservices-template-2024/internal/conf"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	consultants_data "microservices-template-2024/pkg/consultants/data"
	consultants_server "microservices-template-2024/pkg/consultants/server"
	consultants_service "microservices-template-2024/pkg/consultants/service"

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
