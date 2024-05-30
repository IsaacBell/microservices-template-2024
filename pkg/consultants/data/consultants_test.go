package consultants_data_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	consultants_data "microservices-template-2024/pkg/consultants/data"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockDB struct {
	mock.Mock
}

const id = "fake-uuid"

func (m *mockDB) Save(consultant *consultants_biz.Consultant) error {
	args := m.Called(consultant)
	return args.Error(0)
}

func (m *mockDB) First(consultant *consultants_biz.Consultant, id string) error {
	args := m.Called(consultant, id)
	return args.Error(0)
}

func (m *mockDB) Delete(consultant *consultants_biz.Consultant) error {
	args := m.Called(consultant)
	return args.Error(0)
}

func (m *mockDB) Find(consultants *[]*consultants_biz.Consultant) error {
	args := m.Called(consultants)
	return args.Error(0)
}

func TestConsultantRepo_Get(t *testing.T) {
	// mockDB := new(mockDB)
	// logger := log.NewStdLogger(nil)
	testQuery := "SELECT id, user_id FROM `urls` WHERE `id` = $1"
	testDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer testDb.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 testDb,
		PreferSimpleProtocol: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	fmt.Println("here", db)
	rows := sqlmock.NewRows([]string{"id", "user_id"}).
		AddRow("c1", "u1")
	mock.ExpectQuery(regexp.QuoteMeta(testQuery)).
		WillReturnRows(rows).WithArgs(id)

	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, nil)

	// consultant := &consultants_biz.Consultant{ID: id}
	// server.DB.Create(&consultant)
	result, err := repo.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, id)
	assert.Equal(t, result.UserID, "u1")
}

func TestConsultantRepo_Save(t *testing.T) {
	mockDB := new(mockDB)
	logger := log.NewStdLogger(nil)
	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	mockDB.On("Save", consultant).Return(nil).Once()
	result, err := repo.Save(context.Background(), consultant)
	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
	mockDB.AssertExpectations(t)

	consultant = &consultants_biz.Consultant{}
	mockDB.On("Save", consultant).Return(nil).Once()
	result, err = repo.Save(context.Background(), consultant)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	mockDB.AssertExpectations(t)
}

func TestConsultantRepo_Update(t *testing.T) {
	mockDB := new(mockDB)
	logger := log.NewStdLogger(nil)
	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	mockDB.On("Save", consultant).Return(nil).Once()
	result, err := repo.Update(context.Background(), consultant)
	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
	mockDB.AssertExpectations(t)
}

func TestConsultantRepo_Delete(t *testing.T) {
	mockDB := new(mockDB)
	logger := log.NewStdLogger(nil)
	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	mockDB.On("First", consultant, "1").Return(nil).Once()
	mockDB.On("Delete", consultant).Return(nil).Once()
	err := repo.Delete(context.Background(), "1")
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)

	mockDB.On("First", consultant, "2").Return(errors.New("consultant not found")).Once()
	err = repo.Delete(context.Background(), "2")
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestConsultantRepo_Search(t *testing.T) {
	mockDB := new(mockDB)
	logger := log.NewStdLogger(nil)
	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, logger)

	consultant1 := &consultants_biz.Consultant{ID: "1", UserID: "user1", Specializations: []string{"IT"}, Languages: []string{"English"}, YearsOfExperience: 5}
	consultant2 := &consultants_biz.Consultant{ID: "2", UserID: "user2", Specializations: []string{"Marketing"}, Languages: []string{"Spanish"}, YearsOfExperience: 3}
	consultants := []*consultants_biz.Consultant{consultant1, consultant2}
	mockDB.On("Find", &consultants).Return(nil).Once()

	filters := map[string]interface{}{
		"user_id":                 "user1",
		"specializations":         []string{"IT"},
		"languages":               []string{"English"},
		"years_of_experience_gte": 4,
		"years_of_experience_lte": 6,
	}
	result, err := repo.Search(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, consultant1)
	assert.Contains(t, result, consultant2)
	mockDB.AssertExpectations(t)
}

func TestConsultantRepo_SaveCommunication(t *testing.T) {
	logger := log.NewStdLogger(nil)
	data := &consultants_data.Data{}
	repo := consultants_data.NewConsultantRepo(data, logger)

	comm := &consultants_biz.Communication{
		UserID:      "user1",
		RecipientID: "1",
		Msg:         "Test message",
		CommType:    consultants_biz.COMM_TYPE_FromClient,
		Options:     map[string]bool{"option1": true},
		From:        "sender@example.com",
	}
	result, err := repo.SaveCommunication(context.Background(), comm)
	assert.NoError(t, err)
	assert.Equal(t, comm, result)
}
