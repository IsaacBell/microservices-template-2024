package data

import (
	"context"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	if err := server.DB.Where(biz.User{Email: u.Email}).FirstOrCreate(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	if err := server.DB.Updates(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	var u *biz.User
	if err := server.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	var u *biz.User
	if err := server.DB.First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) ListAll(context.Context) ([]*biz.User, error) {
	return nil, nil
}
