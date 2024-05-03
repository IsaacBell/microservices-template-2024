package leads_service

import (
	"github.com/google/wire"
)

// Declare server types to run concurrently at runtime
var ProviderSet = wire.NewSet(
	NewLeadService, NewCompanyService,
)
