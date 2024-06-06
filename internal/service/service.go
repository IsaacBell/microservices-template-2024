package service

import (
	finService "core/pkg/finance/service"

	"github.com/google/wire"
)

// Declare services available at runtime
var ProviderSet = wire.NewSet(
	NewGreeterService, NewUsersService, NewTransactionsService, NewLiabilitiesService, NewLogService,
	finService.NewFinanceService,
)
