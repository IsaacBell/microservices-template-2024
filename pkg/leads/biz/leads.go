package leads_biz

import (
	"context"
	"fmt"
	leadsV1 "microservices-template-2024/api/v1/b2b"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	ID        string `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email     string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Address1  string `protobuf:"bytes,3,opt,name=address1,proto3" json:"address1,omitempty"`
	Address2  string `protobuf:"bytes,4,opt,name=address2,proto3" json:"address2,omitempty"`
	Phone     string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	City      string `protobuf:"bytes,6,opt,name=city,proto3" json:"city,omitempty"`
	State     string `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	Zipcode   string `protobuf:"bytes,8,opt,name=zipcode,proto3" json:"zipcode,omitempty"`
	Country   string `protobuf:"bytes,9,opt,name=country,proto3" json:"country,omitempty"`
	Title     string `protobuf:"bytes,10,opt,name=title,proto3" json:"title,omitempty"`
	CompanyID string `protobuf:"bytes,11,opt,name=company_id,json=companyId,proto3" json:"company_id,omitempty"`
	Company   *Company
	Position  string `protobuf:"bytes,13,opt,name=position,proto3" json:"position,omitempty"`
	Timezone  string `protobuf:"bytes,14,opt,name=timezone,proto3" json:"timezone,omitempty"`
	Owner     string `protobuf:"bytes,15,opt,name=owner,proto3" json:"owner,omitempty"`
	Summary   string `protobuf:"bytes,16,opt,name=summary,proto3" json:"summary,omitempty"`
	Deleted   bool   `protobuf:"varint,17,opt,name=deleted,proto3" json:"deleted,omitempty"`
	Synced    bool   `protobuf:"varint,18,opt,name=synced,proto3" json:"synced,omitempty"`
}

func (t *Lead) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	return nil
}

func LeadToProtoData(lead *Lead) *leadsV1.Lead {
	if lead == nil {
		return nil
	}

	return &leadsV1.Lead{
		Id:        lead.ID,
		Email:     lead.Email,
		Address1:  lead.Address1,
		Address2:  lead.Address2,
		Phone:     lead.Phone,
		City:      lead.City,
		State:     lead.State,
		Zipcode:   lead.Zipcode,
		Country:   lead.Country,
		Title:     lead.Title,
		CompanyId: lead.CompanyID,
		Company:   CompanyToProtoData(lead.Company),
		Position:  lead.Position,
		Timezone:  lead.Timezone,
		Owner:     lead.Owner,
		Summary:   lead.Summary,
		Deleted:   lead.Deleted,
		Synced:    lead.Synced,
	}
}

func ProtoToLeadData(input *leadsV1.Lead) *Lead {
	lead := &Lead{
		ID:        input.Id,
		Email:     input.Email,
		Address1:  input.Address1,
		Address2:  input.Address2,
		Phone:     input.Phone,
		City:      input.City,
		State:     input.State,
		Zipcode:   input.Zipcode,
		Country:   input.Country,
		Title:     input.Title,
		CompanyID: input.CompanyId,
		Company:   ProtoToCompanyData(input.Company),
		Position:  input.Position,
		Timezone:  input.Timezone,
		Owner:     input.Owner,
		Summary:   input.Summary,
		Deleted:   input.Deleted,
		Synced:    input.Synced,
	}

	return lead
}

type LeadRepo interface {
	Get(context.Context, string) (*Lead, error)
	Save(context.Context, *Lead) (*Lead, error)
	Update(context.Context, *Lead) (*Lead, error)
	Delete(context.Context, string) error
}

type LeadAction struct {
	repo LeadRepo
	log  *log.Helper
}

func NewLeadAction(repo LeadRepo, logger log.Logger) *LeadAction {
	return &LeadAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *LeadAction) GetLead(ctx context.Context, id string) (*Lead, error) {
	uc.log.WithContext(ctx).Infof("GetLead: %s", id)
	lead, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return lead, nil
}

func (uc *LeadAction) CreateLead(ctx context.Context, u *Lead) (*Lead, error) {
	uc.log.WithContext(ctx).Infof("CreateLead: %s", u.Email)
	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		fmt.Println("error creating lead: ", err)
	}
	fmt.Println("create lead result: ", res)
	return res, err
}

func (uc *LeadAction) Update(ctx context.Context, u *Lead) (*Lead, error) {
	uc.log.WithContext(ctx).Infof("UpdateLead: %s", u.Email)
	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		fmt.Println("error updating lead: ", err)
	}
	fmt.Println("update lead result: ", res)
	return res, err
}

func (uc *LeadAction) Delete(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("Delete Lead: %s", id)
	return uc.repo.Delete(ctx, id)
}
