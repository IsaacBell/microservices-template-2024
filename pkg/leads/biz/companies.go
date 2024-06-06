package leads_biz

import (
	"context"
	leadsV1 "core/api/v1/b2b"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID          string              `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email       string              `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Address1    string              `protobuf:"bytes,3,opt,name=address1,proto3" json:"address1,omitempty"`
	Address2    string              `protobuf:"bytes,4,opt,name=address2,proto3" json:"address2,omitempty"`
	Phone       string              `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	City        string              `protobuf:"bytes,6,opt,name=city,proto3" json:"city,omitempty"`
	State       string              `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	Zipcode     string              `protobuf:"bytes,8,opt,name=zipcode,proto3" json:"zipcode,omitempty"`
	Country     string              `protobuf:"bytes,9,opt,name=country,proto3" json:"country,omitempty"`
	Name        string              `protobuf:"bytes,10,opt,name=name,proto3" json:"name,omitempty"`
	Domain      string              `protobuf:"bytes,11,opt,name=domain,proto3" json:"domain,omitempty"`
	Industry    string              `protobuf:"bytes,12,opt,name=industry,proto3" json:"industry,omitempty"`
	Description string              `protobuf:"bytes,13,opt,name=description,proto3" json:"description,omitempty"`
	Type        leadsV1.CompanyType `protobuf:"varint,14,opt,name=type,proto3,enum=api.v1.CompanyType" json:"type,omitempty"`
	Deleted     bool                `protobuf:"varint,15,opt,name=deleted,proto3" json:"deleted,omitempty"`
	Synced      bool                `protobuf:"varint,16,opt,name=synced,proto3" json:"synced,omitempty"`
}

func (t *Company) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	return nil
}

func CompanyToProtoData(company *Company) *leadsV1.Company {
	if company == nil {
		return nil
	}

	return &leadsV1.Company{
		Id:          company.ID,
		Email:       company.Email,
		Address1:    company.Address1,
		Address2:    company.Address2,
		Phone:       company.Phone,
		City:        company.City,
		State:       company.State,
		Zipcode:     company.Zipcode,
		Country:     company.Country,
		Name:        company.Name,
		Domain:      company.Domain,
		Industry:    company.Industry,
		Description: company.Description,
		Type:        company.Type,
		Deleted:     company.Deleted,
		Synced:      company.Synced,
	}
}

func ProtoToCompanyData(input *leadsV1.Company) *Company {
	company := &Company{
		ID:          input.Id,
		Email:       input.Email,
		Address1:    input.Address1,
		Address2:    input.Address2,
		Phone:       input.Phone,
		City:        input.City,
		State:       input.State,
		Zipcode:     input.Zipcode,
		Country:     input.Country,
		Name:        input.Name,
		Domain:      input.Domain,
		Industry:    input.Industry,
		Description: input.Description,
		Type:        input.Type,
		Deleted:     input.Deleted,
		Synced:      input.Synced,
	}

	return company
}

type CompanyRepo interface {
	Get(context.Context, string) (*Company, error)
	Save(context.Context, *Company) (*Company, error)
	Update(context.Context, *Company) (*Company, error)
	Delete(context.Context, string) error
}

type CompanyAction struct {
	repo CompanyRepo
	log  *log.Helper
}

func NewCompanyAction(repo CompanyRepo, logger log.Logger) *CompanyAction {
	return &CompanyAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *CompanyAction) GetCompany(ctx context.Context, id string) (*Company, error) {
	uc.log.WithContext(ctx).Infof("GetCompany: %s", id)
	company, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (uc *CompanyAction) CreateCompany(ctx context.Context, u *Company) (*Company, error) {
	uc.log.WithContext(ctx).Infof("CreateCompany: %s", u.Email)
	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		fmt.Println("error creating company: ", err)
	}
	fmt.Println("create company result: ", res)
	return res, err
}

func (uc *CompanyAction) Update(ctx context.Context, u *Company) (*Company, error) {
	uc.log.WithContext(ctx).Infof("UpdateCompany: %s", u.Email)
	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		fmt.Println("error updating company: ", err)
	}
	fmt.Println("update company result: ", res)
	return res, err
}

func (uc *CompanyAction) Delete(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("Delete Company: %s", id)
	return uc.repo.Delete(ctx, id)
}
