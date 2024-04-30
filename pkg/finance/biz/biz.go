package finance_biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewStockQuoteAction,
)
