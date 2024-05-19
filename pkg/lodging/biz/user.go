package lodging_biz

import (
	"context"
	"fmt"

	v1 "microservices-template-2024/api/v1"

	"github.com/go-kratos/kratos/v2/log"
)

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

func UserToProtoData(user *User) *v1.User {
	if user == nil {
		return nil
	}

	return &v1.User{
		Id:           user.ID,
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
	}
}

func ProtoToUserData(input *v1.User) *User {
	user := &User{}
	user.ID = input.Id
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

	// if input.CreatedAt != nil {
	// 	user.CreatedAt = input.CreatedAt
	// }
	// if input.UpdatedAt != nil {
	// 	user.UpdatedAt = input.UpdatedAt
	// }
	// if input.LastLoginAt != nil {
	// 	user.LastLoginAt = input.LastLoginAt
	// }

	user.Timezone = input.Timezone
	user.Locale = input.Locale
	// user.Metadata = input.Metadata
	user.Deleted = input.Deleted

	return user
}

func NewUserAction(repo UserRepo, logger log.Logger) *UserAction {
	return &UserAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserAction) CreateUser(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %s", u.Email)
	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		fmt.Println("error creating user: ", err)
	}
	fmt.Println("create user result: ", res)
	return res, err
}

func (uc *UserAction) UpdateUser(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %s", u.Email)
	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		fmt.Println("error updating user: ", err)
	}
	fmt.Println("update user result: ", res)
	return res, err
}

func (uc *UserAction) FindUserById(ctx context.Context, id string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindUser: %s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *UserAction) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	uc.log.WithContext(ctx).Infof("FindUser: %s", email)
	return uc.repo.FindByEmail(ctx, email)
}

func (uc *UserAction) Delete(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("Delete User: %s", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *UserAction) ListAll(ctx context.Context) ([]*User, error) {
	uc.log.WithContext(ctx).Infof("List Users")
	return uc.repo.ListAll(ctx)
}
