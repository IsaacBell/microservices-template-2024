package data

import (
	"context"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"

	"github.com/go-kratos/kratos/v2/log"
)

type transactionRepo struct {
	data *Data
	log  *log.Helper
}

func NewTransactionRepo(data *Data, logger log.Logger) biz.TransactionRepo {
	return &transactionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *transactionRepo) Save(ctx context.Context, u *biz.Transaction) (*biz.Transaction, error) {
	if u.ID != "" {
		if err := server.DB.Save(&u).Error; err != nil {
			return nil, err
		} else {
			return u, nil
		}
	}

	if err := server.DB.Omit("ID").FirstOrCreate(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *transactionRepo) Update(ctx context.Context, u *biz.Transaction) (*biz.Transaction, error) {
	if err := server.DB.Save(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *transactionRepo) FindByID(ctx context.Context, id string) (*biz.Transaction, error) {
	var u *biz.Transaction
	if err := server.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *transactionRepo) ListAll(context.Context) ([]*biz.Transaction, error) {
	var transactions []*biz.Transaction

	if err := server.DB.Last(&transactions).Limit(100).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
