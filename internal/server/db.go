package server

import (
	"fmt"
	"log"
	"microservices-template-2024/internal/biz"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func automigrateDBTables(*gorm.DB) {
	DB.AutoMigrate(&biz.User{})
	DB.AutoMigrate(&biz.Transaction{})
}

func AmountGT(db *gorm.DB, amt int) *gorm.DB {
	return db.Where("amount > ?", amt)
}

func AmountLT(db *gorm.DB, amt int) *gorm.DB {
	return db.Where("amount < ?", amt)
}

func Unsynced(db *gorm.DB) *gorm.DB {
	return db.Where("synced = ?", false)
}

func Synced(db *gorm.DB) *gorm.DB {
	return db.Where("synced = ?", true)
}

func DbConnString() string {
	user := os.Getenv("COCKROACH_DB_USER")
	pass := os.Getenv("COCKROACH_DB_PASS")
	url := os.Getenv("COCKROACH_DB_URL")
	db := os.Getenv("COCKROACH_DB_DBNAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=verify-full", user, pass, url, db)
	log.Println("dsn: ", dsn)

	return dsn
}

func OpenDBConn() error {
	dsn := DbConnString()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	var now time.Time
	DB.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)

	automigrateDBTables(DB)

	return nil
}
