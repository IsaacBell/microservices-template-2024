package server

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbString() string {
	user := os.Getenv("COCKROACH_DB_USER")
	pass := os.Getenv("COCKROACH_DB_PASS")
	url := os.Getenv("COCKROACH_DB_URL")
	db := os.Getenv("COCKROACH_DB_DBNAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=verify-full", user, pass, url, db)
	log.Println("dsn: ", dsn)

	return dsn
}

func OpenDBConn() {
	dsn := DbString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)
}
