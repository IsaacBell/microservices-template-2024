package data

import (
	"context"
	"fmt"
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"

	"github.com/go-kratos/kratos/v2/log"
)

type liabilityRepo struct {
	data *Data
	log  *log.Helper
}

func NewLiabilityRepo(data *Data, logger log.Logger) biz.LiabilityRepo {
	return &liabilityRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *liabilityRepo) Save(ctx context.Context, l *biz.Liability) (*biz.Liability, error) {
	if l.ID != "" {
		if err := server.DB.Save(&l).Error; err != nil {
			return nil, err
		} else {
			return l, nil
		}
	}

	if err := l.BeforeCreate(server.DB); err != nil {
		return nil, err
	}
	
	if err := server.DB.FirstOrCreate(&l).Error; err != nil {
		return nil, err
	}
	fmt.Println("Liability ID: ", l.ID)

	return l, nil
}

func (r *liabilityRepo) Update(ctx context.Context, u *biz.Liability) (*biz.Liability, error) {
	if err := server.DB.Save(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *liabilityRepo) FindByID(ctx context.Context, id string) (*biz.Liability, error) {
	var l *biz.Liability
	if err := server.DB.First(&l, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (r *liabilityRepo) GetLiabilities(ctx context.Context, req *v1.GetLiabilitiesRequest) ([]*biz.Liability, error) {
	var liabilities []*biz.Liability

	query := server.DB.Limit(300)
	if req.Owner != "" {
		query = query.Where("account_id = ?", req.Owner)
	}
	if req.Since != nil {
		query = query.Where("created_at >= ?", req.Since.AsTime())
	}

	if err := query.Find(&liabilities).Error; err != nil {
		return nil, err
	}

	return liabilities, nil
}

func (r *liabilityRepo) ListAll(context.Context) ([]*biz.Liability, error) {
	var liabilities []*biz.Liability

	if err := server.DB.Last(&liabilities).Limit(100).Error; err != nil {
		return nil, err
	}

	return liabilities, nil
}
