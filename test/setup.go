package test

import (
	"microservices-template-2024/internal/biz"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&biz.User{})
	db.AutoMigrate(&biz.Transaction{})

	return db, nil
}
