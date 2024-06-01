package test

import (
	"microservices-template-2024/internal/server"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	onceOpenDbConn sync.Once
)

func Logger() log.Logger {
	return log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
}

func TestDB() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	onceOpenDbConn.Do(func() {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			server.SetupDbTables(db)
		}
	})
	if err != nil {
		return db, err
	}

	return db, nil
}

func SetupAndTeardownTest(testType string, args ...interface{}) func() {
	// I want to call this like defer SetupAndTeardownTest("consultants_service")()
	// In fact, maybe I want to split into sep. setup & teardown fns
	var onTeardown func()
	switch testType {
	case "":
		onTeardown = defaultTeardown
	case "consultants_service":
		onTeardown = consultantsServiceTeardown
	}
	return func() {
		onTeardown()
	}
}

func defaultSetup()    {}
func defaultTeardown() {}

func consultantsServiceSetup() {
	// repo := new(mocks.MockConsultantRepo)
	// logger := log.With(log.NewStdLogger(os.Stdout),
	// 	"ts", log.DefaultTimestamp,
	// 	"caller", log.DefaultCaller,
	// )
	// action := consultants_biz.NewConsultantAction(repo, logger)
	// consultant := &consultants_biz.Consultant{ID: id, YearsOfExperience: 5}
	// repo.On("Get", context.Background(), id).Return(consultant, nil)

	// service := consultants_service.NewConsultantService(action)
}
func consultantsServiceTeardown() {}
