package leads_biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewLeadAction, NewCompanyAction,
)
