package mocks

import (
	"context"
	"fmt"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ConsultantActionMock struct {
	mock.Mock
	*consultants_biz.ConsultantAction
}

func (m *ConsultantActionMock) SetRepo(r consultants_biz.ConsultantRepo) {
	m.Action.SetRepo(r)
}

func (m *ConsultantActionMock) CreateConsultant(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*consultants_biz.Consultant), nil
}

func (m *ConsultantActionMock) DeleteConsultant(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ConsultantActionMock) GetConsultant(ctx context.Context, id string) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consultants_biz.Consultant), nil
}

func (m *ConsultantActionMock) ListConsultants(ctx context.Context, filters map[string]interface{}) ([]*consultants_biz.Consultant, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*consultants_biz.Consultant), nil
}

func (m *ConsultantActionMock) SendComm(ctx context.Context, c *consultants_biz.Communication) (*consultants_biz.Communication, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*consultants_biz.Communication), nil
}

func (m *ConsultantActionMock) UpdateConsultant(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*consultants_biz.Consultant), nil
}

type MockConsultantRepo struct {
	consultants_biz.ConsultantRepo
	mock.Mock
}

func (m *MockConsultantRepo) Get(ctx context.Context, id string) (*consultants_biz.Consultant, error) {
	fmt.Println("there")
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consultants_biz.Consultant), nil
}

func (m *MockConsultantRepo) Save(ctx context.Context, consultant *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, consultant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consultants_biz.Consultant), nil
}

func (m *MockConsultantRepo) Update(ctx context.Context, consultant *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	args := m.Called(ctx, consultant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consultants_biz.Consultant), nil
}

func (m *MockConsultantRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockConsultantRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*consultants_biz.Consultant, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*consultants_biz.Consultant), nil
}

func (m *MockConsultantRepo) SaveCommunication(ctx context.Context, comm *consultants_biz.Communication) (*consultants_biz.Communication, error) {
	args := m.Called(ctx, comm)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*consultants_biz.Communication), nil
}

type MockDB struct {
	mock.Mock
	*gorm.DB
}

func (m *MockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) First(dest interface{}, id interface{}) *MockDB {
	args := m.Called(dest, id)
	return args.Get(0).(*MockDB)
}

func (m *MockDB) FirstOrCreate(dest interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockDB) Delete(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	arguments := m.Called(query, args)
	return arguments.Get(0).(*MockDB)
}

func (m *MockDB) Find(dest interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}
