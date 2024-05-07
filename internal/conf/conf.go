package conf

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func ConfigDir() string {
	_, err := os.Stat("./configs")
	if os.IsNotExist(err) {
		return "../../configs/"
	}

	return "./configs/"
}

// CSV copied from Plaid API
func CategoriesCsvPath() string {
	return ConfigDir() + "transactions-personal-finance-category-taxonomy.csv"
}

// https://docs.google.com/spreadsheets/d/1I3pBxjfXB056-g_JYf_6o3Rns3BV2kMGG1nCatb91ls/edit?usp=sharing
func FinnhubExchangesCsvPath() string {
	return ConfigDir() + "finnhub-exchanges.csv"
}

func PersonalFinanceCategories() [][]string {
	// Skip the header row
	return readCsvFile(CategoriesCsvPath())[1:]
}

func FinnhubExchanges() [][]string {
	// Skip the header row
	return readCsvFile(FinnhubExchangesCsvPath())[1:]
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

// Redis Caching

func RedisConn(ctx context.Context) *redis.Client {
	pass := os.Getenv("UPSTASH_REDIS_PASS")
	url := os.Getenv("UPSTASH_REDIS_URL")
	path := "rediss://default:" + pass + "@" + url
	opt, _ := redis.ParseURL(path)
	client := redis.NewClient(opt)

	val := client.Get(ctx, "lastConnectedAt").Val()
	fmt.Println("Redis - lastConnectedAt: ", val)
	client.Set(ctx, "lastConnectedAt", time.Now().String(), 0)

	return client
}
