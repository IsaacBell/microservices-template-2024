package consultants_biz

import (
	"context"
	"fmt"
	consultantV1 "microservices-template-2024/api/v1/consultants"
	"microservices-template-2024/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Consultant struct {
	gorm.Model
	ID                string                 `gorm:"primaryKey" protobuf:"bytes,9,opt,name=id,proto3" json:"id,omitempty"`
	UserID            string                 `gorm:"index" protobuf:"bytes,10,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User              *biz.User              `gorm:"foreignKey:UserID" protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Specializations   []string               `protobuf:"bytes,2,rep,name=specializations,proto3" json:"specializations,omitempty"`
	Bio               string                 `protobuf:"bytes,3,opt,name=bio,proto3" json:"bio,omitempty"`
	Languages         []string               `protobuf:"bytes,4,rep,name=languages,proto3" json:"languages,omitempty"`
	YearsOfExperience int32                  `protobuf:"varint,5,opt,name=years_of_experience,json=yearsOfExperience,proto3" json:"years_of_experience,omitempty"`
	Certifications    []string               `protobuf:"bytes,6,rep,name=certifications,proto3" json:"certifications,omitempty"`
	Education         []string               `protobuf:"bytes,7,rep,name=education,proto3" json:"education,omitempty"`
	AdditionalFields  map[string]string      `protobuf:"bytes,8,rep,name=additional_fields,json=additionalFields,proto3" json:"additional_fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CreatedAt         *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt         *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type Comm struct {
	gorm.Model
	Msg         string                         `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	UserID      string                         `gorm:"index" protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CommType    consultantV1.CommunicationType `protobuf:"varint,3,opt,name=comm_type,json=commType,proto3,enum=api.v1.consultants.CommunicationType" json:"comm_type,omitempty"`
	Options     map[string]bool                `gorm:"type:jsonb" protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	RecipientID string                         `protobuf:"bytes,5,opt,name=recipient_id,json=recipientId,proto3" json:"recipient_id,omitempty"`
}

func (c *Consultant) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func ConsultantToProtoData(consultant *Consultant) *consultantV1.Consultant {
	if consultant == nil {
		return nil
	}

	return &consultantV1.Consultant{
		User:              biz.UserToProtoData(consultant.User),
		Specializations:   consultant.Specializations,
		Bio:               consultant.Bio,
		Languages:         consultant.Languages,
		YearsOfExperience: consultant.YearsOfExperience,
		Certifications:    consultant.Certifications,
		Education:         consultant.Education,
		AdditionalFields:  consultant.AdditionalFields,
		Id:                consultant.ID,
		UserId:            consultant.UserID,
	}
}

func ProtoToConsultantData(input *consultantV1.Consultant) *Consultant {
	consultant := &Consultant{
		ID:                input.Id,
		UserID:            input.UserId,
		User:              biz.ProtoToUserData(input.User),
		Specializations:   input.Specializations,
		Bio:               input.Bio,
		Languages:         input.Languages,
		YearsOfExperience: input.YearsOfExperience,
		Certifications:    input.Certifications,
		Education:         input.Education,
		AdditionalFields:  input.AdditionalFields,
		// CreatedAt:         input.CreatedAt,
		// UpdatedAt:         input.UpdatedAt,
	}

	return consultant
}

type ConsultantRepo interface {
	Get(context.Context, string) (*Consultant, error)
	Save(context.Context, *Consultant) (*Consultant, error)
	Update(context.Context, *Consultant) (*Consultant, error)
	Delete(context.Context, string) error
	Search(context.Context, map[string]interface{}) ([]*Consultant, error)
	SaveCommunication(context.Context, *consultantV1.Comm) (*consultantV1.Ack, error)
}

type ConsultantAction struct {
	repo ConsultantRepo
	log  *log.Helper
}

func NewConsultantAction(repo ConsultantRepo, logger log.Logger) *ConsultantAction {
	return &ConsultantAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ConsultantAction) GetConsultant(ctx context.Context, id string) (*Consultant, error) {
	uc.log.WithContext(ctx).Infof("GetConsultant: %s", id)
	consultant, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return consultant, nil
}

func (uc *ConsultantAction) ListConsultants(ctx context.Context, filters map[string]interface{}) ([]*Consultant, error) {
	uc.log.WithContext(ctx).Infof("ListConsultants: %s", filters)
	consultants, err := uc.repo.Search(ctx, filters)
	if err != nil {
		return nil, err
	}

	return consultants, nil
}

func (uc *ConsultantAction) SendComm(ctx context.Context, req *consultantV1.SendCommsRequest) (*consultantV1.SendCommsReply, error) {
	uc.log.WithContext(ctx).Infof("SendComm: %s", req.Comm)
	ack, err := uc.repo.SaveCommunication(ctx, req.Comm)
	if err != nil {
		return nil, err
	}

	return &consultantV1.SendCommsReply{
		Ok:  err == nil,
		Ack: ack,
	}, err
}

func (uc *ConsultantAction) CreateConsultant(ctx context.Context, c *Consultant) (*Consultant, error) {
	uc.log.WithContext(ctx).Infof("CreateConsultant: %s", c.ID)
	res, err := uc.repo.Save(ctx, c)
	if err != nil {
		fmt.Println("error creating consultant: ", err)
	}
	fmt.Println("create consultant result: ", res)
	return res, err
}

func (uc *ConsultantAction) UpdateConsultant(ctx context.Context, c *Consultant) (*Consultant, error) {
	uc.log.WithContext(ctx).Infof("UpdateConsultant: %s", c.ID)
	res, err := uc.repo.Update(ctx, c)
	if err != nil {
		fmt.Println("error updating consultant: ", err)
	}
	fmt.Println("update consultant result: ", res)
	return res, err
}

func (uc *ConsultantAction) DeleteConsultant(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("Delete Consultant: %s", id)
	return uc.repo.Delete(ctx, id)
}
