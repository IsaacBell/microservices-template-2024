package consultants_biz_test

import (
	"context"
	"errors"
	"testing"

	communicationsV1 "microservices-template-2024/api/v1/communications"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConsultantRepo struct {
	mock.Mock
}

func (m *MockConsultantRepo) Get(ctx context.Context, id string) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*consultants_biz.Consultant), args.Error(1)
}

func (m *MockConsultantRepo) Save(ctx context.Context, consultant *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, consultant)
	return args.Get(0).(*consultants_biz.Consultant), args.Error(1)
}

func (m *MockConsultantRepo) Update(ctx context.Context, consultant *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, consultant)
	return args.Get(0).(*consultants_biz.Consultant), args.Error(1)
}

func (m *MockConsultantRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockConsultantRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*consultants_biz.Consultant, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*consultants_biz.Consultant), args.Error(1)
}

func (m *MockConsultantRepo) SaveCommunication(ctx context.Context, comm *consultants_biz.Communication) (*consultants_biz.Communication, error) {
	args := m.Called(ctx, comm)
	return args.Get(0).(*consultants_biz.Communication), args.Error(1)
}

func TestGetConsultant(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}

	result, err := action.GetConsultant(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestListConsultants(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultantList := []*consultants_biz.Consultant{{ID: "1"}, {ID: "2"}}
	filters := map[string]interface{}{"specialization": "IT"}
	repo.On("Search", mock.Anything, filters).Return(consultantList, nil)

	result, err := action.ListConsultants(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, consultantList, result)
}

func TestSendComm(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	comm := &consultants_biz.Communication{UserID: "1", CommType: consultants_biz.CommunicationType(communicationsV1.CommunicationType_from_client.String())}
	repo.On("SaveCommunication", mock.Anything, comm).Return(comm, nil)

	result, err := action.SendComm(context.Background(), comm)
	assert.NoError(t, err)
	assert.Equal(t, comm, result)
}

func TestCreateConsultant(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	repo.On("Save", mock.Anything, consultant).Return(consultant, nil)

	result, err := action.CreateConsultant(context.Background(), consultant)
	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestUpdateConsultant(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	repo.On("Update", mock.Anything, consultant).Return(consultant, nil)

	result, err := action.UpdateConsultant(context.Background(), consultant)
	assert.NoError(t, err)
	assert.Equal(t, consultant, result)
}

func TestDeleteConsultant(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	repo.On("Delete", mock.Anything, "1").Return(nil)

	err := action.DeleteConsultant(context.Background(), "1")
	assert.NoError(t, err)
}

func TestGetConsultantError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	repo.On("Get", mock.Anything, "1").Return(nil, errors.New("consultant not found"))

	_, err := action.GetConsultant(context.Background(), "1")
	assert.Error(t, err)
}

func TestListConsultantsError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	filters := map[string]interface{}{"specialization": "IT"}
	repo.On("Search", mock.Anything, filters).Return(nil, errors.New("error searching consultants"))

	_, err := action.ListConsultants(context.Background(), filters)
	assert.Error(t, err)
}

func TestSendCommError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	comm := &consultants_biz.Communication{UserID: "1", CommType: consultants_biz.CommunicationType(communicationsV1.CommunicationType_from_consultant.String())}
	repo.On("SaveCommunication", mock.Anything, comm).Return(nil, errors.New("error sending communication"))

	_, err := action.SendComm(context.Background(), comm)
	assert.Error(t, err)
}

func TestCreateConsultantError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	repo.On("Save", mock.Anything, consultant).Return(nil, errors.New("error creating consultant"))

	_, err := action.CreateConsultant(context.Background(), consultant)
	assert.Error(t, err)
}

func TestUpdateConsultantError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	consultant := &consultants_biz.Consultant{ID: "1"}
	repo.On("Update", mock.Anything, consultant).Return(nil, errors.New("error updating consultant"))

	_, err := action.UpdateConsultant(context.Background(), consultant)
	assert.Error(t, err)
}

func TestDeleteConsultantError(t *testing.T) {
	repo := new(MockConsultantRepo)
	logger := log.NewStdLogger(nil)
	action := consultants_biz.NewConsultantAction(repo, logger)

	repo.On("Delete", mock.Anything, "1").Return(errors.New("error deleting consultant"))

	err := action.DeleteConsultant(context.Background(), "1")
	assert.Error(t, err)
}
