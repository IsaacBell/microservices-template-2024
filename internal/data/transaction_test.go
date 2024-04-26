package data_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/data"
	"microservices-template-2024/internal/server"
	"microservices-template-2024/test"
)

func TestTransactionRepo(t *testing.T) {
	// Set up the test database
	db, err := test.SetupTestDB()
	require.NoError(t, err)

	server.DB = db

	repo := data.NewTransactionRepo(nil, nil)

	t.Run("Save", func(t *testing.T) {
		transaction := &biz.Transaction{
			ID: "1",
		}

		result, err := repo.Save(context.Background(), transaction)
		assert.NoError(t, err)
		assert.Equal(t, transaction.ID, result.ID)
	})

	t.Run("Update", func(t *testing.T) {
		transaction := &biz.Transaction{
			ID: "1",
		}

		result, err := repo.Update(context.Background(), transaction)
		assert.NoError(t, err)
		assert.Equal(t, transaction.ID, result.ID)
	})

	t.Run("FindByID", func(t *testing.T) {
		result, err := repo.FindByID(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "1", result.ID)
	})

	t.Run("ListAll", func(t *testing.T) {
		result, err := repo.ListAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}
