package finance_biz

import (
	"context"
	v1 "microservices-template-2024/api/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type StockQuote struct {
	gorm.Model
	ID            string  `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Price         float64 `protobuf:"bytes,2,opt,name=price,proto3" json:"price,omitempty"`
	Change        float64 `protobuf:"bytes,3,opt,name=change,proto3" json:"change,omitempty"`
	PercentChange float64 `protobuf:"bytes,4,opt,name=percent_change,proto3" json:"percent_change,omitempty"`
	High          float64 `protobuf:"bytes,5,opt,name=high,proto3" json:"high,omitempty"`
	Low           float64 `protobuf:"bytes,6,opt,name=low,proto3" json:"low,omitempty"`
	Open          float64 `protobuf:"bytes,7,opt,name=open,proto3" json:"open,omitempty"`
	LastClose     float64 `protobuf:"bytes,8,opt,name=last_close,proto3" json:"last_close,omitempty"`
	Symbol        string  `protobuf:"bytes,9,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Synced        bool    `protobuf:"name=synced,proto3" json:"synced,omitempty"`
	Deleted       bool    `protobuf:"name=deleted,proto3" json:"deleted,omitempty"`
}

func FinnHubQuoteToDBRecord(q *finnhub.Quote, symbol string) *StockQuote {
	return &StockQuote{
		Price:         float64(*q.C),
		PercentChange: float64(*q.Dp),
		High:          float64(*q.H),
		Low:           float64(*q.L),
		Open:          float64(*q.O),
		LastClose:     float64(*q.Pc),
		Symbol:        symbol,
		Synced:        false,
		Deleted:       false,
	}
}

func (t *StockQuote) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	return nil
}

func QuoteToProtoData(quote *StockQuote) *v1.StockQuote {
	if quote == nil {
		return nil
	}

	return &v1.StockQuote{
		Price:         quote.Price,
		Change:        quote.Change,
		PercentChange: quote.PercentChange,
		High:          quote.High,
		Low:           quote.Low,
		Open:          quote.Open,
		LastClose:     quote.LastClose,
		Symbol:        quote.Symbol,
	}
}

func ProtoToQuoteData(input *v1.StockQuote) *StockQuote {
	quote := &StockQuote{
		Price:         input.Price,
		Change:        input.Change,
		PercentChange: input.PercentChange,
		High:          input.High,
		Low:           input.Low,
		Open:          input.Open,
		LastClose:     input.LastClose,
		Symbol:        input.Symbol,
		Synced:        false,
		Deleted:       false,
	}

	return quote
}

type StockQuoteRepo interface {
	Get(context.Context, string) (*StockQuote, error)
	WatchTrades(cowner string, stream v1.Finance_WatchTradesServer) error
}

type StockQuoteAction struct {
	repo StockQuoteRepo
	log  *log.Helper
}

func NewStockQuoteAction(repo StockQuoteRepo, logger log.Logger) *StockQuoteAction {
	return &StockQuoteAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *StockQuoteAction) GetStockQuote(ctx context.Context, symbol string) (*v1.GetStockQuoteReply, error) {
	uc.log.WithContext(ctx).Infof("GetStockQuote: %s", symbol)
	quote, err := uc.repo.Get(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return &v1.GetStockQuoteReply{
		Symbol: quote.Symbol,
		Quote:  QuoteToProtoData(quote),
	}, nil
}

func (uc *StockQuoteAction) WatchTrades(symbol string, stream v1.Finance_WatchTradesServer) error {
	uc.log.WithContext(stream.Context()).Infof("WatchTrades")
	return uc.repo.WatchTrades(symbol, stream)
}
