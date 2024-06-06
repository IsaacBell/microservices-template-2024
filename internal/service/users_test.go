package service_test

import (
	"context"
	"errors"
	"testing"

	v1 "core/api/v1"
	"core/internal/biz"
	"core/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserAction struct {
	mock.Mock
	biz.UserAction
}

func (m *mockUserAction) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*biz.User), args.Error(1)
}

func (m *mockUserAction) UpdateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*biz.User), args.Error(1)
}

func (m *mockUserAction) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserAction) FindUserById(ctx context.Context, id string) (*biz.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*biz.User), args.Error(1)
}

func (m *mockUserAction) FindUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*biz.User), args.Error(1)
}

func (m *mockUserAction) ListAll(ctx context.Context) ([]*biz.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*biz.User), args.Error(1)
}

type userActionMock struct {
	action biz.UserAction
	mock   *mockUserAction
	repo   *biz.UserRepo
	mock.Mock
	biz.UserAction
}

func (u *userActionMock) SetMock(mock *mockUserAction) {
	u.mock = mock
	// u.UserAction = mock
}

func TestUsersService_CreateUser(t *testing.T) {
	mockAction := &mockUserAction{}
	svc := service.NewUsersService(&mockAction.UserAction)

	user := &biz.User{
		Email:    "test@example.com",
		Username: "testuser",
	}
	mockAction.On("CreateUser", mock.Anything, user).Return(user, nil)

	req := &v1.CreateUserRequest{User: biz.UserToProtoData(user)}
	reply, err := svc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, user.ID, reply.Id)
}

func TestUsersService_UpdateUser(t *testing.T) {
	mockAction := &mockUserAction{}
	svc := service.NewUsersService(&mockAction.UserAction)

	user := &biz.User{
		ID:       "1",
		Email:    "test@example.com",
		Username: "testuser",
	}
	mockAction.On("UpdateUser", mock.Anything, user).Return(user, nil)

	req := &v1.UpdateUserRequest{User: biz.UserToProtoData(user)}
	reply, err := svc.UpdateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, user.ID, reply.Id)
}

func TestUsersService_DeleteUser(t *testing.T) {
	mockAction := &mockUserAction{}
	svc := service.NewUsersService(&mockAction.UserAction)

	userID := "1"
	mockAction.On("Delete", mock.Anything, userID).Return(nil)

	req := &v1.DeleteUserRequest{Id: userID}
	reply, err := svc.DeleteUser(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
}

func TestUsersService_GetUser(t *testing.T) {
	mockAction := &mockUserAction{}
	svc := service.NewUsersService(&mockAction.UserAction)

	user := &biz.User{
		ID:       "1",
		Email:    "test@example.com",
		Username: "testuser",
	}

	t.Run("get user by ID", func(t *testing.T) {
		mockAction.On("FindUserById", mock.Anything, user.ID).Return(user, nil)

		req := &v1.GetUserRequest{Id: &user.ID}
		reply, err := svc.GetUser(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, biz.UserToProtoData(user), reply.User)
	})

	t.Run("get user by email", func(t *testing.T) {
		mockAction.On("FindUserByEmail", mock.Anything, user.Email).Return(user, nil)

		req := &v1.GetUserRequest{Email: &user.Email}
		reply, err := svc.GetUser(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, biz.UserToProtoData(user), reply.User)
	})

	t.Run("user not found", func(t *testing.T) {
		mockAction.On("FindUserById", mock.Anything, "2").Return(nil, errors.New("user not found"))

		req := &v1.GetUserRequest{Id: &[]string{"2"}[0]}
		reply, err := svc.GetUser(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, reply)
	})
}

func TestUsersService_ListUser(t *testing.T) {
	mockAction := &mockUserAction{}
	svc := service.NewUsersService(&mockAction.UserAction)

	users := []*biz.User{
		{ID: "1", Email: "user1@example.com", Username: "user1"},
		{ID: "2", Email: "user2@example.com", Username: "user2"},
	}
	mockAction.On("ListAll", mock.Anything).Return(users, nil)

	req := &v1.ListUserRequest{}
	reply, err := svc.ListUser(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, reply.Ok)
	assert.Equal(t, len(users), len(reply.Users))
	for i, user := range users {
		assert.Equal(t, biz.UserToProtoData(user), reply.Users[i])
	}
}
