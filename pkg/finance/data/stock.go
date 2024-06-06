package finance_data

import (
	"context"
	v1 "core/api/v1"
	"core/internal/server"
	biz "core/pkg/finance/biz"
	finance_util "core/pkg/finance/util"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type stockQuoteRepo struct {
	data *Data
	log  *log.Helper
}

func NewStockQuoteRepo(data *Data, logger log.Logger) biz.StockQuoteRepo {
	return &stockQuoteRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *stockQuoteRepo) Get(ctx context.Context, symbol string) (*biz.StockQuote, error) {
	var quote *biz.StockQuote
	res, _, err := finance_util.GetFinnhubClient().Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}

	quote = biz.FinnHubQuoteToDBRecord(&res, symbol)
	if quote == nil {
		return nil, err
	}

	if err := server.DB.Create(&quote).Error; err != nil {
		log.Debug("Error saving stock quote to DB")
	}

	return quote, nil
}

func (r *stockQuoteRepo) WatchTrades(symbol string, stream v1.Finance_WatchTradesServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var quotes []*biz.StockQuote
			query := server.DB.Limit(300).Where("synced != ?", true)
			if symbol != "" {
				query = query.Where("symbol = ?", symbol)
			}
			query = query.Find(&quotes)
			if err := query.Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// No unsynced quotes found, wait for a short duration and continue the loop
					time.Sleep(1 * time.Second)
					continue
				}

				return err
			}

			for _, quote := range quotes {
				fmt.Println("sync stock quote: ", quote.Symbol)
				if !quote.Synced {
					quote.Synced = true
					if err := server.DB.Save(&quote).Error; err != nil {
						return err
					}
				}

				if err := stream.Send(&v1.SyncTradesReply{}); err != nil {
					return err
				}
			}

			// Loop time
			time.Sleep(400 * time.Millisecond)
		}
	}

	return nil
}
