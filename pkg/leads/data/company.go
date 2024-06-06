package leads_data

import (
	"context"
	"microservices-template-2024/internal/server"
	leads_biz "microservices-template-2024/pkg/leads/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type companyRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyRepo(data *Data, logger log.Logger) leads_biz.CompanyRepo {
	return &companyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *companyRepo) Get(ctx context.Context, id string) (*leads_biz.Company, error) {
	var company *leads_biz.Company
	err := server.DB.Scopes(server.Active).First(&company, id).Error
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (r *companyRepo) Save(ctx context.Context, u *leads_biz.Company) (*leads_biz.Company, error) {
	if u.ID != "" {
		if err := server.DB.Save(&u).Error; err != nil {
			return nil, err
		} else {
			return u, nil
		}
	}

	if err := server.DB.FirstOrCreate(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *companyRepo) Update(ctx context.Context, company *leads_biz.Company) (*leads_biz.Company, error) {
	if err := server.DB.Save(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepo) Delete(ctx context.Context, id string) error {
	var company *leads_biz.Company
	if err := server.DB.Scopes(server.Active).First(&company, id).Error; err != nil {
		return err
	}
	company.Deleted = true

	if err := server.DB.Save(&company).Error; err != nil {
		return err
	}
	return nil
}
