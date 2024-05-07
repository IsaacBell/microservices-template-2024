package service

import (
	finService "microservices-template-2024/pkg/finance/service"

	"github.com/google/wire"
)

// Declare services available at runtime
var ProviderSet = wire.NewSet(
	NewGreeterService, NewUsersService, NewTransactionsService, NewLiabilitiesService, NewLogService,
	finService.NewFinanceService,
)
