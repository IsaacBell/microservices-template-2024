package service

import "github.com/google/wire"

// Declare services available at runtime
var ProviderSet = wire.NewSet(NewGreeterService, NewUsersService, NewTransactionsService, NewLiabilitiesService)
