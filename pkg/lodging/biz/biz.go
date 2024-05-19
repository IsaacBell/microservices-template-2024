package lodging_biz

import (
	"microservices-template-2024/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewPropertyAction, NewUserAction,
)

type userRepo struct {
	data *conf.Data
	log  *log.Helper
}
