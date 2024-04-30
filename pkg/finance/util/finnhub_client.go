package finance_util

import (
	"os"
	"sync"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var (
	client *finnhub.DefaultApiService
	once   sync.Once
)

func InitFinnhubClient() {
	once.Do(func() {
		cfg := finnhub.NewConfiguration()
		cfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv("FINNHUB_API_TOKEN"))
		client = finnhub.NewAPIClient(cfg).DefaultApi
	})
}

func GetFinnhubClient() *finnhub.DefaultApiService {
	return client
}
