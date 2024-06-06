package data_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"core/internal/biz"
	"core/internal/data"
	"core/internal/server"
	"core/test"
)

func TestUserRepo(t *testing.T) {
	// Set up the test database
	db, err := test.SetupTestDB()
	require.NoError(t, err)

	server.DB = db

	repo := data.NewUserRepo(nil, nil)

	t.Run("Save", func(t *testing.T) {
		user := &biz.User{
			ID:    "1",
			Email: "test@example.com",
		}

		result, err := repo.Save(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("Update", func(t *testing.T) {
		user := &biz.User{
			ID:    "1",
			Email: "updated@example.com",
		}

		result, err := repo.Update(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("FindByID", func(t *testing.T) {
		result, err := repo.FindByID(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "updated@example.com", result.Email)
	})

	t.Run("FindByEmail", func(t *testing.T) {
		result, err := repo.FindByEmail(context.Background(), "updated@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "updated@example.com", result.Email)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(context.Background(), "1")
		assert.NoError(t, err)

		_, err = repo.FindByID(context.Background(), "1")
		assert.Error(t, err)
	})

	t.Run("ListAll", func(t *testing.T) {
		user := &biz.User{
			Email: "list@example.com",
		}
		_, err := repo.Save(context.Background(), user)
		assert.NoError(t, err)

		result, err := repo.ListAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}
