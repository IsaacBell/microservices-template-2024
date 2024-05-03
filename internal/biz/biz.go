package biz

import (
	"time"

	"github.com/google/wire"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var ProviderSet = wire.NewSet(
	NewGreeterUsecase, NewUserAction, NewTransactionAction, NewLiabilityAction,
)

func TimestampToProto(ts time.Time) *timestamppb.Timestamp {
	return timestamppb.New(ts)
}
