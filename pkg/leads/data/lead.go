package leads_data

import (
	"context"
	"microservices-template-2024/internal/server"
	leads_biz "microservices-template-2024/pkg/leads/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type leadRepo struct {
	data *Data
	log  *log.Helper
}

func NewLeadRepo(data *Data, logger log.Logger) leads_biz.LeadRepo {
	return &leadRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *leadRepo) Get(ctx context.Context, id string) (*leads_biz.Lead, error) {
	var lead *leads_biz.Lead
	err := server.DB.Scopes(server.Active).First(&lead, id).Error
	if err != nil {
		return nil, err
	}

	return lead, nil
}

func (r *leadRepo) Save(ctx context.Context, u *leads_biz.Lead) (*leads_biz.Lead, error) {
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

func (r *leadRepo) Update(ctx context.Context, lead *leads_biz.Lead) (*leads_biz.Lead, error) {
	if err := server.DB.Save(&lead).Error; err != nil {
		return nil, err
	}
	return lead, nil
}

func (r *leadRepo) Delete(ctx context.Context, id string) error {
	var lead *leads_biz.Lead
	if err := server.DB.Scopes(server.Active).First(&lead, id).Error; err != nil {
		return err
	}
	lead.Deleted = true

	if err := server.DB.Save(&lead).Error; err != nil {
		return err
	}
	return nil
}
