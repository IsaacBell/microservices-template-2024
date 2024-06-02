package data

import (
	"context"
	"fmt"
	"time"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"
	"microservices-template-2024/internal/util"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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

func (r *transactionRepo) Save(ctx context.Context, t *biz.Transaction) (*biz.Transaction, error) {
	defer util.Benchmark("transactionRepo.Save()")()

	if t.ID != "" {
		if err := server.DB.Save(&t).Error; err != nil {
			return nil, err
		} else {
			return t, nil
		}
	}

	if err := t.BeforeCreate(server.DB); err != nil {
		return nil, err
	}

	if err := server.DB.FirstOrCreate(&t).Error; err != nil {
		return nil, err
	}
	fmt.Println("Transaction ID: ", t.ID)

	return t, nil
}

func (r *transactionRepo) Update(ctx context.Context, u *biz.Transaction) (*biz.Transaction, error) {
	defer util.Benchmark("transactionRepo.Update()")()
	if err := server.DB.Save(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *transactionRepo) FindByID(ctx context.Context, id string) (*biz.Transaction, error) {
	defer util.Benchmark("transactionRepo.FindByID()")()
	var u *biz.Transaction
	if err := server.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *transactionRepo) SyncTransactions(ctx context.Context, owner string, stream v1.Transactions_SyncTransactionsServer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var transactions []*biz.Transaction
			query := server.DB.Limit(300).Where("synced != ?", true)
			if owner != "" {
				query = query.Where("account_id = ?", owner)
			}
			query = query.First(&transactions)
			if err := query.Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// No unsynced transactions found, wait for a short duration and continue the loop
					time.Sleep(1 * time.Second)
					continue
				}

				return err
			}

			for _, transaction := range transactions {
				fmt.Println("sync transaction: ", transaction.ID)
				if !transaction.Synced {
					transaction.Synced = true
					if err := server.DB.Save(&transaction).Error; err != nil {
						return err
					}
				}

				if err := stream.Send(&v1.GetTransactionsReply{Transaction: biz.TransactionToProtoData(transaction)}); err != nil {
					return err
				}
			}

			// Loop time
			time.Sleep(150 * time.Millisecond)
		}
	}

	return nil
}

func (r *transactionRepo) ListAll(context.Context) ([]*biz.Transaction, error) {
	defer util.Benchmark("transactionRepo.ListAll()")()
	var transactions []*biz.Transaction

	if err := server.DB.Last(&transactions).Limit(100).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
