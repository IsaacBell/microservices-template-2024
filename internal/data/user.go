package data

import (
	"context"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"
	"microservices-template-2024/internal/util"

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
	defer util.Benchmark("userRepo.Save()")()

	// if err := server.DB.Where(biz.User{Email: u.Email}).FirstOrCreate(&u).Error; err != nil {
	if u.ID != "" {
		if err := server.DB.Save(&u).Error; err != nil {
			return nil, err
		} else {
			return u, nil
		}
	}

	u.Deleted = false

	if err := server.DB.Omit("ID").FirstOrCreate(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	defer util.Benchmark("userRepo.Update()")()

	if err := server.DB.Save(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	defer util.Benchmark("userRepo.FindByID()")()
	var u *biz.User
	if err := server.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	defer util.Benchmark("userRepo.FindByEmail()")()
	var u *biz.User
	if err := server.DB.First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	defer util.Benchmark("userRepo.Delete()")()

	if err := server.DB.Delete(&biz.User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ListAll(context.Context) ([]*biz.User, error) {
	defer util.Benchmark("userRepo.ListAll()")()

	var users []*biz.User

	if err := server.DB.Last(&users).Limit(100).Error; err != nil {
		return nil, err
	}

	return users, nil
}
