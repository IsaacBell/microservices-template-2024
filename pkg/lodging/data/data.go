package lodging_data

import (
	"core/internal/conf"
	"core/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewData, NewPropertyRepo, data.NewUserRepo,
)

// Data .
type Data struct {
	// wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{}, cleanup, nil
}
