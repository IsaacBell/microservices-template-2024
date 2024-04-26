package test

import (
	"microservices-template-2024/internal/biz"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Logger() log.Logger {
	return log.NewStdLogger(os.Stdout)
}

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&biz.User{})
	db.AutoMigrate(&biz.Transaction{})

	return db, nil
}
