package biz

import (
	"context"
	"fmt"
	v1 "microservices-template-2024/api/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID                             string                      `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccountID                      string                      `protobuf:"bytes,2,opt,name=account_id,proto3" json:"account_id,omitempty"`
	AccountOwner                   string                      `protobuf:"bytes,3,opt,name=account_owner,proto3" json:"account_owner,omitempty"`
	Amount                         float64                     `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	IsoCurrencyCode                string                      `protobuf:"bytes,5,opt,name=iso_currency_code,proto3" json:"iso_currency_code,omitempty"`
	UnofficialCurrencyCode         string                      `protobuf:"bytes,6,opt,name=unofficial_currency_code,proto3" json:"unofficial_currency_code,omitempty"`
	Category                       []string                    `gorm:"type:text[]" protobuf:"bytes,7,rep,name=category,proto3" json:"category,omitempty"`
	CategoryID                     string                      `protobuf:"bytes,8,opt,name=category_id,proto3" json:"category_id,omitempty"`
	CheckNumber                    string                      `protobuf:"bytes,9,opt,name=check_number,proto3" json:"check_number,omitempty"`
	Counterparties                 []*v1.Counterparty          `gorm:"type:jsonb" protobuf:"bytes,10,rep,name=counterparties,proto3" json:"counterparties,omitempty"`
	Date                           string                      `protobuf:"bytes,11,opt,name=date,proto3" json:"date,omitempty"`
	Datetime                       string                      `protobuf:"bytes,12,opt,name=datetime,proto3" json:"datetime,omitempty"`
	AuthorizedDate                 string                      `protobuf:"bytes,13,opt,name=authorized_date,proto3" json:"authorized_date,omitempty"`
	AuthorizedDatetime             string                      `protobuf:"bytes,14,opt,name=authorized_datetime,proto3" json:"authorized_datetime,omitempty"`
	Location                       *v1.Location                `gorm:"embedded" protobuf:"bytes,15,opt,name=location,proto3" json:"location,omitempty"`
	Name                           string                      `protobuf:"bytes,16,opt,name=name,proto3" json:"name,omitempty"`
	MerchantName                   string                      `protobuf:"bytes,17,opt,name=merchant_name,proto3" json:"merchant_name,omitempty"`
	MerchantEntityID               string                      `protobuf:"bytes,18,opt,name=merchant_entity_id,proto3" json:"merchant_entity_id,omitempty"`
	LogoURL                        string                      `protobuf:"bytes,19,opt,name=logo_url,proto3" json:"logo_url,omitempty"`
	Website                        string                      `protobuf:"bytes,20,opt,name=website,proto3" json:"website,omitempty"`
	PaymentMeta                    *v1.PaymentMeta             `gorm:"embedded" protobuf:"bytes,21,opt,name=payment_meta,proto3" json:"payment_meta,omitempty"`
	PaymentChannel                 string                      `protobuf:"bytes,22,opt,name=payment_channel,proto3" json:"payment_channel,omitempty"`
	Pending                        bool                        `protobuf:"varint,23,opt,name=pending,proto3" json:"pending,omitempty"`
	PendingTransactionID           string                      `protobuf:"bytes,24,opt,name=pending_transaction_id,proto3" json:"pending_transaction_id,omitempty"`
	PersonalFinanceCategory        *v1.PersonalFinanceCategory `gorm:"embedded" protobuf:"bytes,25,opt,name=personal_finance_category,proto3" json:"personal_finance_category,omitempty"`
	PersonalFinanceCategoryIconURL string                      `protobuf:"bytes,26,opt,name=personal_finance_category_icon_url,proto3" json:"personal_finance_category_icon_url,omitempty"`
	TransactionID                  string                      `protobuf:"bytes,27,opt,name=transaction_id,proto3" json:"transaction_id,omitempty"`
	TransactionCode                string                      `protobuf:"bytes,28,opt,name=transaction_code,proto3" json:"transaction_code,omitempty"`
	TransactionType                string                      `protobuf:"bytes,29,opt,name=transaction_type,proto3" json:"transaction_type,omitempty"`
	Synced                         bool                        `protobuf:"bytes,30,opt,name=synced,proto3" json:"synced,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (*Transaction, error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	return t, nil
}

func TransactionToProtoData(trans *Transaction) *v1.Transaction {
	if trans == nil {
		return nil
	}

	return &v1.Transaction{
		Id:                             trans.ID,
		AccountId:                      trans.AccountID,
		AccountOwner:                   trans.AccountOwner,
		Amount:                         trans.Amount,
		IsoCurrencyCode:                trans.IsoCurrencyCode,
		UnofficialCurrencyCode:         trans.UnofficialCurrencyCode,
		Category:                       trans.Category,
		CategoryId:                     trans.CategoryID,
		CheckNumber:                    trans.CheckNumber,
		Counterparties:                 trans.Counterparties,
		Date:                           trans.Date,
		Datetime:                       trans.Datetime,
		AuthorizedDate:                 trans.AuthorizedDate,
		AuthorizedDatetime:             trans.AuthorizedDatetime,
		Location:                       trans.Location,
		Name:                           trans.Name,
		MerchantName:                   trans.MerchantName,
		MerchantEntityId:               trans.MerchantEntityID,
		LogoUrl:                        trans.LogoURL,
		Website:                        trans.Website,
		PaymentMeta:                    trans.PaymentMeta,
		PaymentChannel:                 trans.PaymentChannel,
		Pending:                        trans.Pending,
		PendingTransactionId:           trans.PendingTransactionID,
		PersonalFinanceCategory:        trans.PersonalFinanceCategory,
		PersonalFinanceCategoryIconUrl: trans.PersonalFinanceCategoryIconURL,
		Synced:                         trans.Synced,
		TransactionId:                  trans.TransactionID,
		TransactionCode:                trans.TransactionCode,
		TransactionType:                trans.TransactionType,
	}
}

func ProtoToTransactionData(input *v1.Transaction) *Transaction {
	trans := &Transaction{}
	trans.ID = input.Id
	trans.AccountID = input.AccountId
	trans.AccountOwner = input.AccountOwner
	trans.Amount = input.Amount
	trans.IsoCurrencyCode = input.IsoCurrencyCode
	trans.UnofficialCurrencyCode = input.UnofficialCurrencyCode
	trans.Category = input.Category
	trans.CategoryID = input.CategoryId
	trans.CheckNumber = input.CheckNumber
	trans.Counterparties = input.Counterparties
	trans.Date = input.Date
	trans.Datetime = input.Datetime
	trans.AuthorizedDate = input.AuthorizedDate
	trans.AuthorizedDatetime = input.AuthorizedDatetime
	trans.Location = input.Location
	trans.Name = input.Name
	trans.MerchantName = input.MerchantName
	trans.MerchantEntityID = input.MerchantEntityId
	trans.LogoURL = input.LogoUrl
	trans.Website = input.Website
	trans.PaymentMeta = input.PaymentMeta
	trans.PaymentChannel = input.PaymentChannel
	trans.Pending = input.Pending
	trans.PendingTransactionID = input.PendingTransactionId
	trans.PersonalFinanceCategory = input.PersonalFinanceCategory
	trans.PersonalFinanceCategoryIconURL = input.PersonalFinanceCategoryIconUrl
	trans.Synced = input.Synced
	trans.TransactionID = input.TransactionId
	trans.TransactionCode = input.TransactionCode
	trans.TransactionType = input.TransactionType

	return trans
}

type TransactionRepo interface {
	Save(context.Context, *Transaction) (*Transaction, error)
	Update(context.Context, *Transaction) (*Transaction, error)
	FindByID(context.Context, string) (*Transaction, error)
	SyncTransactions(ctx context.Context, owner string, stream v1.Transactions_SyncTransactionsServer) error
	ListAll(context.Context) ([]*Transaction, error)
}

type TransactionAction struct {
	repo TransactionRepo
	log  *log.Helper
}

func NewTransactionAction(repo TransactionRepo, logger log.Logger) *TransactionAction {
	return &TransactionAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *TransactionAction) CreateTransaction(ctx context.Context, t *Transaction) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("CreateTransaction: %s", t.TransactionID)
	res, err := uc.repo.Save(ctx, t)
	if err != nil {
		fmt.Println("error creating transaction: ", err)
	}
	fmt.Println("create transaction result: ", res)
	return res, err
}

// *fmt.wrapError {msg: "unsupported data type: &[]", err: error(*errors.errorString) ...}

func (uc *TransactionAction) UpdateTransaction(ctx context.Context, t *Transaction) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("UpdateTransaction: %s", t.TransactionID)
	res, err := uc.repo.Update(ctx, t)
	if err != nil {
		fmt.Println("error updating transaction: ", err)
	}
	fmt.Println("update transaction result: ", res)
	return res, err
}

func (uc *TransactionAction) FindTransactionByID(ctx context.Context, id string) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("FindTransaction: %s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *TransactionAction) SyncTransactions(ctx context.Context, owner string, stream v1.Transactions_SyncTransactionsServer) error {
	uc.log.WithContext(ctx).Infof("Sync Transactions")
	return uc.repo.SyncTransactions(ctx, owner, stream)
}

func (uc *TransactionAction) ListAll(ctx context.Context) ([]*Transaction, error) {
	uc.log.WithContext(ctx).Infof("List Transactions")
	return uc.repo.ListAll(ctx)
}
