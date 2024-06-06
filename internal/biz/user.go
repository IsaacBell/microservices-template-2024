package biz

import (
	"context"
	v1 "core/api/v1"
	"fmt"
	"time"

	// v1 "core/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	moesifModels "github.com/moesif/moesifapi-go/models"
)

// var (
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

type User struct {
	gorm.Model
	ID           string `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username     string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email        string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	PasswordHash string `protobuf:"bytes,4,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	FirstName    string `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	PhoneNumber  string `protobuf:"bytes,7,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	AvatarURL    string `protobuf:"bytes,8,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	// Roles         []string `protobuf:"bytes,9,rep,name=roles,proto3" json:"roles,omitempty"`
	EmailVerified bool `protobuf:"varint,10,opt,name=email_verified,json=emailVerified,proto3" json:"email_verified,omitempty"`
	PhoneVerified bool `protobuf:"varint,11,opt,name=phone_verified,json=phoneVerified,proto3" json:"phone_verified,omitempty"`
	// CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// LastLoginAt   *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=last_login_at,json=lastLoginAt,proto3" json:"last_login_at,omitempty"`
	Timezone string `protobuf:"bytes,15,opt,name=timezone,proto3" json:"timezone,omitempty"`
	Locale   string `protobuf:"bytes,16,opt,name=locale,proto3" json:"locale,omitempty"`
	// Metadata map[string]string `protobuf:"bytes,17,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Deleted      bool   `protobuf:"bytes,19,name=deleted,proto3" json:"deleted,omitempty"`
	Title        string `protobuf:"bytes,23,name=title,proto3" json:"title,proto3" json:"title,omitempty"`
	CompanyId    string `protobuf:"bytes,24,name=company_id,proto3" json:"company_id,proto3" json:"company_id,omitempty"`
	SessionToken string `protobuf:"bytes,25,name=session_token,proto3" json:"session_token,proto3" json:"session_token,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func UserToProtoData(user *User) *v1.User {
	if user == nil {
		return nil
	}

	return &v1.User{
		Id:           user.ID,
		CompanyId:    user.CompanyId,
		SessionToken: user.SessionToken,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
		AvatarUrl:    user.AvatarURL,
		// Roles:         user.Roles,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		// CreatedAt:     user.CreatedAt,
		// UpdatedAt:     user.UpdatedAt,
		// LastLoginAt:   user.LastLoginAt,
		Timezone: user.Timezone,
		Locale:   user.Locale,
		// Metadata: user.Metadata,
		Deleted: user.Deleted,
		Title:   user.Title,
	}
}

func ProtoToUserData(input *v1.User) *User {
	if input == nil {
		return nil
	}

	user := &User{}
	user.ID = input.Id
	user.CompanyId = input.CompanyId
	user.SessionToken = input.SessionToken
	user.Username = input.Username
	user.Email = input.Email
	user.PasswordHash = input.PasswordHash
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.PhoneNumber = input.PhoneNumber
	user.AvatarURL = input.AvatarUrl
	// user.Roles = input.Roles
	user.EmailVerified = input.EmailVerified
	user.PhoneVerified = input.PhoneVerified

	user.Timezone = input.Timezone
	user.Locale = input.Locale
	// user.Metadata = input.Metadata
	user.Deleted = input.Deleted
	user.Title = input.Title

	return user
}

func UserToMoesifData(u *User) *moesifModels.UserModel {
	userLastModifiedTime := time.Now().Local()
	return &moesifModels.UserModel{
		UserId:       u.ID,
		ModifiedTime: &userLastModifiedTime,
		SessionToken: &u.SessionToken,
		Metadata: map[string]interface{}{
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"title":      u.Title,
			"sales_info": map[string]interface{}{
				"stage":          "Customer",
				"lifetime_value": 0,
				"account_owner":  u.Email,
			},
		},
	}
}

type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, string) error
	FindByID(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
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
	// uc.log.WithContext(ctx).Infof("CreateUser: %s", u.Email)
	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		fmt.Println("error creating user: ", err)
	}
	fmt.Println("create user result: ", res)
	return res, err
}

func (uc *UserAction) UpdateUser(ctx context.Context, u *User) (*User, error) {
	// uc.log.WithContext(ctx).Infof("UpdateUser: %s", u.Email)
	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		fmt.Println("error updating user: ", err)
	}
	fmt.Println("update user result: ", res)
	return res, err
}

func (uc *UserAction) FindUserById(ctx context.Context, id string) (*User, error) {
	// uc.log.WithContext(ctx).Infof("FindUser: %s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *UserAction) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	// uc.log.WithContext(ctx).Infof("FindUser: %s", email)
	return uc.repo.FindByEmail(ctx, email)
}

func (uc *UserAction) Delete(ctx context.Context, id string) error {
	// uc.log.WithContext(ctx).Infof("Delete User: %s", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *UserAction) ListAll(ctx context.Context) ([]*User, error) {
	// uc.log.WithContext(ctx).Infof("List Users")
	return uc.repo.ListAll(ctx)
}
