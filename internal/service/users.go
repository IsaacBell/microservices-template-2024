package service

import (
	"context"
	"errors"
	"fmt"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	cache "microservices-template-2024/pkg/cache"
	stream "microservices-template-2024/pkg/stream"
)

type UsersService struct {
	v1.UnimplementedUsersServer

	action *biz.UserAction
}

func NewUsersService(action *biz.UserAction) *UsersService {
	return &UsersService{action: action}
}

func (s *UsersService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	if req.User == nil {
		return &v1.CreateUserReply{Ok: false, Id: ""}, errors.New("user data not supplied")
	}

	user := biz.ProtoToUserData(req.User)
	fmt.Println(user.Email)
	res, err := s.action.CreateUser(ctx, user)

	fmt.Println("CreateUser response: ", res)
	stream.ProduceKafkaMessage("main", "New User: "+user.Email)

	cache.Cache(ctx).Set("user:"+req.User.Id, user, 0)

	return &v1.CreateUserReply{Ok: err == nil, Id: res.ID}, err
}

func (s *UsersService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	if req.User == nil {
		return &v1.UpdateUserReply{Ok: false, Id: ""}, errors.New("user data not supplied")
	}

	user := biz.ProtoToUserData(req.User)
	res, err := s.action.UpdateUser(ctx, user)
	cache.Cache(ctx).Set("user:"+req.User.Id, user, 0)
	return &v1.UpdateUserReply{Ok: err == nil, Id: res.ID}, err
}

func (s *UsersService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	if req.Id == "" {
		return &v1.DeleteUserReply{Ok: false}, errors.New("id not supplied")
	}

	err := s.action.Delete(ctx, req.Id)
	if err != nil {
		return &v1.DeleteUserReply{Ok: false}, err
	}

	go func() {
		err := cache.Cache(ctx).Del(req.Id).Err()
		if err != nil {
			fmt.Printf("Failed to delete cache entry for user %s: %v \n", req.Id, err)
		}
	}()

	return &v1.DeleteUserReply{Ok: err == nil}, err
}

func (s *UsersService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserReply, error) {
	if *req.Id == "" && *req.Email == "" && *req.Name == "" {
		return &v1.GetUserReply{User: nil}, errors.New("user data not supplied")
	}

	var u *biz.User
	var err error

	if req.Id != nil {
		u, err = s.action.FindUserById(ctx, *req.Id)
	} else if req.Email != nil {
		u, err = s.action.FindUserByEmail(ctx, *req.Email)
	}

	if err != nil {
		return nil, err
	}

	if req.Id != nil {
		cache.Cache(ctx).Set("user+id:"+*req.Id, u, 0)
	} else {
		cache.Cache(ctx).Set("user+email:"+*req.Email, u, 0)
	}

	return &v1.GetUserReply{User: biz.UserToProtoData(u)}, nil
}

func (s *UsersService) ListUser(ctx context.Context, req *v1.ListUserRequest) (*v1.ListUserReply, error) {
	list, err := s.action.ListAll(ctx)
	users := make([]*v1.User, len(list))

	for i, u := range list {
		users[i] = biz.UserToProtoData(u)
	}

	return &v1.ListUserReply{Users: users, Ok: err == nil}, err
}
