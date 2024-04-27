// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/conf"
	"microservices-template-2024/internal/data"
	"microservices-template-2024/internal/server"
	"microservices-template-2024/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	userRepo := data.NewUserRepo(dataData, logger)
	userAction := biz.NewUserAction(userRepo, logger)
	usersService := service.NewUsersService(userAction)
	transactionRepo := data.NewTransactionRepo(dataData, logger)
	transactionAction := biz.NewTransactionAction(transactionRepo, logger)
	transactionsService := service.NewTransactionsService(transactionAction)
	liabilityRepo := data.NewLiabilityRepo(dataData, logger)
	liabilityAction := biz.NewLiabilityAction(liabilityRepo, logger)
	liabilitiesService := service.NewLiabilitiesService(liabilityAction)
	grpcServer := server.NewGRPCServer(confServer, greeterService, usersService, transactionsService, liabilitiesService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, usersService, transactionsService, liabilitiesService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
