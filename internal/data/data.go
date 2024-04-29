package data

import (
	"microservices-template-2024/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// Declare data repos available at runtime
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo, NewTransactionRepo, NewLiabilityRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
