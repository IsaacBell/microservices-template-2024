package server

import (
	"fmt"
	"log"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/conf"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	"microservices-template-2024/pkg/notifications"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func automigrateDBTables(*gorm.DB) {
	DB.AutoMigrate(&consultants_biz.Consultant{})
	DB.AutoMigrate(&notifications.Notification{})
	DB.AutoMigrate(&biz.User{})
	DB.AutoMigrate(&biz.Transaction{})
	DB.AutoMigrate(&Exchange{})
	DB.AutoMigrate(&Category{})
	DB.AutoMigrate(&Category{})
}

func AmountGT(db *gorm.DB, amt int) *gorm.DB {
	return db.Where("amount > ?", amt)
}

func AmountLT(db *gorm.DB, amt int) *gorm.DB {
	return db.Where("amount < ?", amt)
}

func Active(db *gorm.DB) *gorm.DB {
	return db.Where("deleted = ?", false)
}

func Unsynced(db *gorm.DB) *gorm.DB {
	return db.Where("synced = ?", false)
}

func Synced(db *gorm.DB) *gorm.DB {
	return db.Where("synced = ?", true)
}

type Category struct {
	gorm.Model
	Primary     string
	Detailed    string
	Description string
}

type Exchange struct {
	gorm.Model
	Code        string
	Name        string
	Mic         string
	Timezone    string
	Premarket   string
	Hour        string
	Postmarket  string
	CloseDate   string
	Country     string
	CountryName string
	Source      string
}

func SeedCategories(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Exchange{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	records := conf.PersonalFinanceCategories()
	categories := make([]*Category, len(records))

	for i, record := range records {
		category := Category{
			Primary:     record[0],
			Detailed:    record[1],
			Description: record[2],
		}
		categories[i] = &category
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	return nil
}

func SeedExchanges(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Exchange{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	records := conf.FinnhubExchanges()
	n := len(records)
	exchanges := make([]*Exchange, 0, n)

	for _, exc := range records {
		exchanges = append(exchanges, &Exchange{
			Code:        exc[0],
			Name:        exc[1],
			Mic:         exc[2],
			Timezone:    exc[3],
			Premarket:   exc[4],
			Hour:        exc[5],
			Postmarket:  exc[6],
			CloseDate:   exc[7],
			Country:     exc[8],
			CountryName: exc[9],
			Source:      exc[10],
		})
	}

	if err := db.Create(&exchanges).Error; err != nil {
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

func SetupDbTables(db *gorm.DB) {
	automigrateDBTables(db)
	if err := SeedCategories(db); err != nil {
		log.Fatalln("Error seeding category data: ", err)
	}
	if err := SeedExchanges(db); err != nil {
		log.Fatalln("Error seeding exchange data: ", err)
	}
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

	SetupDbTables(DB)

	return nil
}
