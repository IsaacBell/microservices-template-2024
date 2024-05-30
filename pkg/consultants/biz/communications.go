package consultants_biz

import (
	"errors"
	"fmt"
	consultantsV1 "microservices-template-2024/api/v1/consultants"

	"gorm.io/gorm"
)

type Communication struct {
	gorm.Model
	Msg         string            `gorm:"not null" protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	UserID      string            `gorm:"index;not null" protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CommType    CommunicationType `gorm:"not null" protobuf:"varint,3,opt,name=comm_type,json=commType,proto3,enum=api.v1.consultants.CommunicationType" json:"comm_type,omitempty"`
	Options     map[string]bool   `gorm:"type:jsonb" protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	RecipientID string            `gorm:"not null" protobuf:"bytes,5,opt,name=recipient_id,json=recipientId,proto3" json:"recipient_id,omitempty"`
	From        string            `gorm:"not null" protobuf:"bytes,6,opt,name=from,json=from,proto3" json:"from,omitempty"`
}

type CommunicationType string

const (
	COMM_TYPE_Unknown    CommunicationType = "unknown"
	COMM_TYPE_FromClient CommunicationType = "from_client"
	COMM_TYPE_FromAdmin  CommunicationType = "from_admin"
	COMM_TYPE_FromSystem CommunicationType = "from_system"
)

var (
	CommunicationTypes map[string]CommunicationType = map[string]CommunicationType{
		COMM_TYPE_Unknown.String():    COMM_TYPE_Unknown,
		COMM_TYPE_FromClient.String(): COMM_TYPE_FromClient,
		COMM_TYPE_FromAdmin.String():  COMM_TYPE_FromAdmin,
		COMM_TYPE_FromSystem.String(): COMM_TYPE_FromSystem,
	}
)

func (ct CommunicationType) String() string {
	if ct == "" {
		return COMM_TYPE_Unknown.String()
	}
	return string(ct)
}

func (ct CommunicationType) ToProto() consultantsV1.CommunicationType {
	switch ct {
	case COMM_TYPE_FromClient:
		return consultantsV1.CommunicationType_from_client
	case COMM_TYPE_FromAdmin:
		return consultantsV1.CommunicationType_from_admin
	case COMM_TYPE_FromSystem:
		return consultantsV1.CommunicationType_from_system
	default:
		return consultantsV1.CommunicationType_unknown
	}
}

func FromString(str string) (CommunicationType, error) {
	if ct, ok := CommunicationTypes[str]; ok {
		return ct, nil
	}
	return COMM_TYPE_Unknown, errors.New("unknown communication type: " + str)
}

func CommunicationTypeFromProto(ct consultantsV1.CommunicationType) CommunicationType {
	fmt.Println("Proto message type: ", ct.String())
	return CommunicationType(ct.String())
}
