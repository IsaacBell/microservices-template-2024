package biz

import (
	"context"
	"fmt"
	v1 "microservices-template-2024/api/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Liability struct {
	gorm.Model

	ID        string    `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccountID string    `protobuf:"bytes,2,opt,name=account_id,proto3" json:"account_id,omitempty"`
	Aprs      []*v1.Apr `gorm:"foreignKey:CreditLiabilityID" protobuf:"bytes,2,rep,name=aprs,proto3" json:"aprs,omitempty"`

	CreditLiability   *CreditLiability      `gorm:"embedded" protobuf:"bytes,3,opt,name=credit_liability,proto3" json:"credit_liability,omitempty"`
	MortgageLiability *v1.MortgageLiability `gorm:"embedded" protobuf:"bytes,4,opt,name=mortgage_liability,proto3" json:"mortgage_liability,omitempty"`
	StudentLiability  *v1.StudentLiability  `gorm:"embedded" protobuf:"bytes,5,opt,name=student_liability,proto3" json:"student_liability,omitempty"`
}

type CreditLiability struct {
	v1.CreditLiability
	ID                     string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Aprs                   []*v1.Apr `gorm:"foreignKey:CreditLiabilityID" protobuf:"bytes,2,rep,name=aprs,proto3" json:"aprs,omitempty"`
	IsOverdue              bool      `protobuf:"varint,3,opt,name=is_overdue,proto3" json:"is_overdue,omitempty"`
	LastPaymentAmount      float64   `protobuf:"fixed64,4,opt,name=last_payment_amount,proto3" json:"last_payment_amount,omitempty"`
	LastPaymentDate        string    `protobuf:"bytes,5,opt,name=last_payment_date,proto3" json:"last_payment_date,omitempty"`
	LastStatementIssueDate string    `protobuf:"bytes,6,opt,name=last_statement_issue_date,proto3" json:"last_statement_issue_date,omitempty"`
	LastStatementBalance   float64   `protobuf:"fixed64,7,opt,name=last_statement_balance,proto3" json:"last_statement_balance,omitempty"`
	MinimumPaymentAmount   float64   `protobuf:"fixed64,8,opt,name=minimum_payment_amount,proto3" json:"minimum_payment_amount,omitempty"`
	NextPaymentDueDate     string    `protobuf:"bytes,9,opt,name=next_payment_due_date,proto3" json:"next_payment_due_date,omitempty"`
	AccountID              string    `protobuf:"bytes,10,opt,name=account_id,proto3" json:"account_id,omitempty"`
}

type Apr struct {
	ID                   string  `gorm:"primaryKey" protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
	AprPercentage        float64 `protobuf:"fixed64,1,opt,name=apr_percentage,proto3" json:"apr_percentage,omitempty"`
	AprType              string  `protobuf:"bytes,2,opt,name=apr_type,proto3" json:"apr_type,omitempty"`
	BalanceSubjectToApr  float64 `protobuf:"fixed64,3,opt,name=balance_subject_to_apr,proto3" json:"balance_subject_to_apr,omitempty"`
	InterestChargeAmount float64 `protobuf:"fixed64,4,opt,name=interest_charge_amount,proto3" json:"interest_charge_amount,omitempty"`
	CreditLiabilityID    string  `gorm:"index"`
}

func (l *Liability) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}

	return nil
}

func LiabilityToProtoData(liability *Liability) *v1.Liability {
	if liability == nil {
		return nil
	}
	result := &v1.Liability{
		Id: liability.ID,
	}

	if liability.CreditLiability != nil {
		result.Liability = &v1.Liability_Credit{Credit: CreditLiabilityToProto(liability.CreditLiability)}
	} else if liability.MortgageLiability != nil {
		result.Liability = &v1.Liability_Mortgage{Mortgage: liability.MortgageLiability}
	} else if liability.StudentLiability != nil {
		result.Liability = &v1.Liability_Student{Student: liability.StudentLiability}
	}

	return result
}

func ProtoToCreditLiability(credit *v1.CreditLiability) *CreditLiability {
	return &CreditLiability{
		ID:                     credit.Id,
		AccountID:              credit.AccountId,
		Aprs:                   credit.Aprs,
		IsOverdue:              credit.IsOverdue,
		LastPaymentAmount:      float64(credit.LastPaymentAmount),
		LastStatementIssueDate: credit.LastStatementIssueDate,
		LastStatementBalance:   float64(credit.LastStatementBalance),
		LastPaymentDate:        credit.LastPaymentDate,
		MinimumPaymentAmount:   float64(credit.MinimumPaymentAmount),
		NextPaymentDueDate:     credit.NextPaymentDueDate,
	}
}

func CreditLiabilityToProto(credit *CreditLiability) *v1.CreditLiability {
	return &v1.CreditLiability{
		Id:                     credit.ID,
		AccountId:              credit.AccountId,
		Aprs:                   credit.Aprs,
		IsOverdue:              credit.IsOverdue,
		LastPaymentAmount:      float32(credit.LastPaymentAmount),
		LastStatementIssueDate: credit.LastStatementIssueDate,
		LastStatementBalance:   float32(credit.LastStatementBalance),
		LastPaymentDate:        credit.LastPaymentDate,
		MinimumPaymentAmount:   float32(credit.MinimumPaymentAmount),
		NextPaymentDueDate:     credit.NextPaymentDueDate,
	}
}

func ProtoToLiabilityData(input *v1.GetLiabilityReply) *Liability {
	liability := &Liability{}
	liability.ID = input.Id

	if input.Liability != nil {
		switch x := input.Liability.Liability.(type) {
		case *v1.Liability_Credit:
			tmp := input.Liability.GetCredit()
			credit := ProtoToCreditLiability(tmp)
			liability.CreditLiability = credit
		case *v1.Liability_Mortgage:
			liability.MortgageLiability = x.Mortgage
		case *v1.Liability_Student:
			liability.StudentLiability = x.Student
		}
	}

	return liability
}

type LiabilityRepo interface {
	Save(context.Context, *Liability) (*Liability, error)
	Update(context.Context, *Liability) (*Liability, error)
	FindByID(context.Context, string) (*Liability, error)
	GetLiabilities(context.Context, *v1.GetLiabilitiesRequest) ([]*Liability, error)
	ListAll(context.Context) ([]*Liability, error)
}

type LiabilityAction struct {
	repo LiabilityRepo
	log  *log.Helper
}

func NewLiabilityAction(repo LiabilityRepo, logger log.Logger) *LiabilityAction {
	return &LiabilityAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *LiabilityAction) CreateLiability(ctx context.Context, l *Liability) (*Liability, error) {
	uc.log.WithContext(ctx).Infof("CreateLiability: %s", l.ID)
	res, err := uc.repo.Save(ctx, l)
	if err != nil {
		fmt.Println("error creating liability: ", err)
	}
	fmt.Println("create liability result: ", res)
	return res, err
}

func (uc *LiabilityAction) UpdateLiability(ctx context.Context, l *Liability) (*Liability, error) {
	uc.log.WithContext(ctx).Infof("UpdateLiability: %s", l.ID)
	res, err := uc.repo.Update(ctx, l)
	if err != nil {
		fmt.Println("error updating liability: ", err)
	}
	fmt.Println("update liability result: ", res)
	return res, err
}

func (uc *LiabilityAction) FindLiabilityByID(ctx context.Context, id string) (*Liability, error) {
	uc.log.WithContext(ctx).Infof("FindLiability: %s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *LiabilityAction) GetLiabilities(ctx context.Context, req *v1.GetLiabilitiesRequest) ([]*Liability, error) {
	uc.log.WithContext(ctx).Infof("GetLiabilities")
	return uc.repo.GetLiabilities(ctx, req)
}

func (uc *LiabilityAction) ListAll(ctx context.Context) ([]*Liability, error) {
	uc.log.WithContext(ctx).Infof("List Liabilities")
	return uc.repo.ListAll(ctx)
}
