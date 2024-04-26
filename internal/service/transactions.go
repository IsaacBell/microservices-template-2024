package service

import (
	"context"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	// log "microservices-template-2024/internal/service/log"
)

type TransactionsService struct {
	v1.UnimplementedTransactionsServer

	action *biz.TransactionAction
}

func NewTransactionsService(action *biz.TransactionAction) *TransactionsService {
	return &TransactionsService{action: action}
}

func (s *TransactionsService) CreateTransaction(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionReply, error) {
	transaction := biz.ProtoToTransactionData(req.Transaction)
	t, err := s.action.CreateTransaction(ctx, transaction)
	return &v1.CreateTransactionReply{Ok: err == nil, Id: t.ID}, err
}

func (s *TransactionsService) UpdateTransaction(ctx context.Context, req *v1.UpdateTransactionsRequest) (*v1.UpdateTransactionsReply, error) {
	transaction := biz.ProtoToTransactionData(req.Transaction)
	_, err := s.action.UpdateTransaction(ctx, transaction)
	return &v1.UpdateTransactionsReply{Ok: err == nil}, err
}

func (s *TransactionsService) GetTransaction(ctx context.Context, req *v1.GetTransactionsRequest) (*v1.GetTransactionsReply, error) {
	var u *biz.Transaction
	var err error
	if req.Id != "" {
		u, err = s.action.FindTransactionByID(ctx, *&req.Id)
	}
	if err != nil {
		return nil, err
	}
	return &v1.GetTransactionsReply{Transaction: biz.TransactionToProtoData(u)}, nil
}

func (s *TransactionsService) SyncTransactions(req *v1.ListTransactionsRequest, stream v1.Transactions_SyncTransactionsServer) error {
	ctx := stream.Context()

	err := s.action.SyncTransactions(ctx, req.Owner, stream)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (s *TransactionsService) ListTransaction(ctx context.Context, req *v1.ListTransactionsRequest) (*v1.ListTransactionsReply, error) {
	list, err := s.action.ListAll(ctx)
	transactions := make([]*v1.Transaction, len(list))

	for i, u := range list {
		transactions[i] = biz.TransactionToProtoData(u)
	}

	return &v1.ListTransactionsReply{Transactions: transactions}, err
}
