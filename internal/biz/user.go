package biz

import (
	"context"

	v1 "microservices-template-2024/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type User struct {
	gorm.Model
	ID            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	PasswordHash  string                 `protobuf:"bytes,4,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	FirstName     string                 `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName      string                 `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	PhoneNumber   string                 `protobuf:"bytes,7,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	AvatarURL     string                 `protobuf:"bytes,8,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	Roles         []string               `protobuf:"bytes,9,rep,name=roles,proto3" json:"roles,omitempty"`
	EmailVerified bool                   `protobuf:"varint,10,opt,name=email_verified,json=emailVerified,proto3" json:"email_verified,omitempty"`
	PhoneVerified bool                   `protobuf:"varint,11,opt,name=phone_verified,json=phoneVerified,proto3" json:"phone_verified,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	LastLoginAt   *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=last_login_at,json=lastLoginAt,proto3" json:"last_login_at,omitempty"`
	Timezone      string                 `protobuf:"bytes,15,opt,name=timezone,proto3" json:"timezone,omitempty"`
	Locale        string                 `protobuf:"bytes,16,opt,name=locale,proto3" json:"locale,omitempty"`
	Metadata      map[string]string      `protobuf:"bytes,17,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) ([]*User, error)
	ListAll(context.Context) ([]*User, error)
}

type UserAction struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserAction(repo UserRepo, logger log.Logger) *UserAction {
	return &UserAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserAction) CreateUser(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u.Email)
	return uc.repo.Save(ctx, u)
}
