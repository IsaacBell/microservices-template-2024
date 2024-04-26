package biz_test

import (
	"context"
	"errors"
	"testing"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/test"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockTransactionRepo struct {
	transactions []*biz.Transaction
}

func (m *mockTransactionRepo) Save(ctx context.Context, t *biz.Transaction) (*biz.Transaction, error) {
	m.transactions = append(m.transactions, t)
	return t, nil
}

func (m *mockTransactionRepo) Update(ctx context.Context, t *biz.Transaction) (*biz.Transaction, error) {
	for i, transaction := range m.transactions {
		if transaction.ID == t.ID {
			m.transactions[i] = t
			return t, nil
		}
	}
	return nil, errors.New("transaction not found")
}

func (m *mockTransactionRepo) FindByID(ctx context.Context, id string) (*biz.Transaction, error) {
	for _, transaction := range m.transactions {
		if transaction.ID == id {
			return transaction, nil
		}
	}
	return nil, errors.New("transaction not found")
}

func (m *mockTransactionRepo) SyncTransactions(ctx context.Context, owner string, stream v1.Transactions_SyncTransactionsServer) error {
	for _, transaction := range m.transactions {
		if transaction.AccountOwner == owner {
			if err := stream.Send(&v1.GetTransactionsReply{Transaction: biz.TransactionToProtoData(transaction)}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *mockTransactionRepo) ListAll(ctx context.Context) ([]*biz.Transaction, error) {
	return m.transactions, nil
}

func TestTransactionAction_CreateTransaction(t *testing.T) {
	repo := &mockTransactionRepo{}
	action := biz.NewTransactionAction(repo, test.Logger())

	transaction := &biz.Transaction{
		TransactionID: "123",
		AccountOwner:  "user1",
	}

	createdTransaction, err := action.CreateTransaction(context.Background(), transaction)
	assert.NoError(t, err)
	assert.Equal(t, transaction.TransactionID, createdTransaction.TransactionID)
	assert.Equal(t, transaction.AccountOwner, createdTransaction.AccountOwner)
}

func TestTransactionAction_UpdateTransaction(t *testing.T) {
	repo := &mockTransactionRepo{
		transactions: []*biz.Transaction{
			{ID: "1", TransactionID: "123", AccountOwner: "user1"},
		},
	}
	action := biz.NewTransactionAction(repo, test.Logger())

	updatedTransaction := &biz.Transaction{
		ID:            "1",
		TransactionID: "123",
		AccountOwner:  "user2",
	}

	result, err := action.UpdateTransaction(context.Background(), updatedTransaction)
	assert.NoError(t, err)
	assert.Equal(t, updatedTransaction.AccountOwner, result.AccountOwner)
}

func TestTransactionAction_FindTransactionByID(t *testing.T) {
	transaction := &biz.Transaction{ID: "1", TransactionID: "123", AccountOwner: "user1"}
	repo := &mockTransactionRepo{transactions: []*biz.Transaction{transaction}}
	action := biz.NewTransactionAction(repo, test.Logger())

	foundTransaction, err := action.FindTransactionByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, foundTransaction.ID)
	assert.Equal(t, transaction.TransactionID, foundTransaction.TransactionID)
	assert.Equal(t, transaction.AccountOwner, foundTransaction.AccountOwner)
}

func TestTransactionAction_SyncTransactions(t *testing.T) {
	transactions := []*biz.Transaction{
		{ID: "1", TransactionID: "123", AccountOwner: "user1"},
		{ID: "2", TransactionID: "456", AccountOwner: "user2"},
	}
	repo := &mockTransactionRepo{transactions: transactions}
	action := biz.NewTransactionAction(repo, test.Logger())

	stream := &mockTransactions_SyncTransactionsServer{}
	err := action.SyncTransactions(context.Background(), "user1", stream)
	assert.NoError(t, err)
	assert.Len(t, stream.Replies, 1)
	assert.Equal(t, transactions[0].TransactionID, stream.Replies[0].Transaction.TransactionId)
}

func TestTransactionAction_ListAll(t *testing.T) {
	transactions := []*biz.Transaction{
		{ID: "1", TransactionID: "123", AccountOwner: "user1"},
		{ID: "2", TransactionID: "456", AccountOwner: "user2"},
	}
	repo := &mockTransactionRepo{transactions: transactions}
	action := biz.NewTransactionAction(repo, test.Logger())

	result, err := action.ListAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
}

type mockTransactions_SyncTransactionsServer struct {
	Replies []*v1.GetTransactionsReply
	Errors  []error
	grpc.ServerStream
}

func (m *mockTransactions_SyncTransactionsServer) Send(reply *v1.GetTransactionsReply) error {
	m.Replies = append(m.Replies, reply)
	return nil
}

func (m *mockTransactions_SyncTransactionsServer) Context() context.Context {
	return context.Background()
}

func (m *mockTransactions_SyncTransactionsServer) RecvMsg(interface{}) error {
	return nil
}
