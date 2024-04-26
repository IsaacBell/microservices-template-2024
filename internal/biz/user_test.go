package biz_test

import (
	"context"
	"errors"
	"testing"

	"microservices-template-2024/internal/biz"
	"microservices-template-2024/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	users []*biz.User
}

func (m *mockUserRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	m.users = append(m.users, u)
	return u, nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	for i, user := range m.users {
		if user.ID == u.ID {
			m.users[i] = u
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	for i, user := range m.users {
		if user.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) ListAll(ctx context.Context) ([]*biz.User, error) {
	return m.users, nil
}

func TestUserAction_CreateUser(t *testing.T) {
	repo := &mockUserRepo{}
	action := biz.NewUserAction(repo, test.Logger())

	user := &biz.User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	createdUser, err := action.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Username, createdUser.Username)
}

func TestUserAction_UpdateUser(t *testing.T) {
	repo := &mockUserRepo{
		users: []*biz.User{
			{ID: "1", Email: "test@example.com", Username: "testuser"},
		},
	}
	action := biz.NewUserAction(repo, test.Logger())

	updatedUser := &biz.User{
		ID:       "1",
		Email:    "updated@example.com",
		Username: "updateduser",
	}

	result, err := action.UpdateUser(context.Background(), updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Email, result.Email)
	assert.Equal(t, updatedUser.Username, result.Username)
}

func TestUserAction_FindUserByID(t *testing.T) {
	user := &biz.User{ID: "1", Email: "test@example.com", Username: "testuser"}
	repo := &mockUserRepo{users: []*biz.User{user}}
	action := biz.NewUserAction(repo, test.Logger())

	foundUser, err := action.FindUserById(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserAction_FindUserByEmail(t *testing.T) {
	user := &biz.User{ID: "1", Email: "test@example.com", Username: "testuser"}
	repo := &mockUserRepo{users: []*biz.User{user}}
	action := biz.NewUserAction(repo, test.Logger())

	foundUser, err := action.FindUserByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserAction_Delete(t *testing.T) {
	user := &biz.User{ID: "1", Email: "test@example.com", Username: "testuser"}
	repo := &mockUserRepo{users: []*biz.User{user}}
	action := biz.NewUserAction(repo, test.Logger())

	err := action.Delete(context.Background(), "1")
	assert.NoError(t, err)
	assert.Empty(t, repo.users)
}

func TestUserAction_ListAll(t *testing.T) {
	users := []*biz.User{
		{ID: "1", Email: "test1@example.com", Username: "testuser1"},
		{ID: "2", Email: "test2@example.com", Username: "testuser2"},
	}
	repo := &mockUserRepo{users: users}
	action := biz.NewUserAction(repo, test.Logger())

	result, err := action.ListAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, users, result)
}
