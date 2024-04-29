package server

import (
	"fmt"
	"log"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/conf"
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

type Category struct {
	Primary     string
	Detailed    string
	Description string
}

func SeedCategories(db *gorm.DB) error {
	records := conf.PersonalFinanceCategories()
	categories := make([]*Category, len(records))

	for _, record := range records {
		category := Category{
			Primary:     record[0],
			Detailed:    record[1],
			Description: record[2],
		}
		categories = append(categories, &category)
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	return nil
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
	if err := SeedCategories(DB); err != nil {
		fmt.Println("Error seeding category data: ", err)
	}

	return nil
}
