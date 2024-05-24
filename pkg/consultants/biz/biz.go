package consultants_biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewConsultantAction,
)
