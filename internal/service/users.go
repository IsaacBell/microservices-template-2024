package service

import (
	"context"
	"fmt"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	// log "microservices-template-2024/internal/service/log"
)

type UsersService struct {
	v1.UnimplementedUsersServer

	action *biz.UserAction
}

func NewUsersService(action *biz.UserAction) *UsersService {
	return &UsersService{action: action}
}

func (s *UsersService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	user := biz.ProtoToUserData(req.User)
	fmt.Println(user.Email)
	res, err := s.action.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CreateUser response: ", res)
	return &v1.CreateUserReply{Ok: true}, nil
}
func (s *UsersService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	return &v1.UpdateUserReply{}, nil
}
func (s *UsersService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	return &v1.DeleteUserReply{Ok: true}, nil
}
func (s *UsersService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserReply, error) {
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
	return &v1.GetUserReply{User: biz.UserToProtoData(u)}, nil
}
func (s *UsersService) ListUser(ctx context.Context, req *v1.ListUserRequest) (*v1.ListUserReply, error) {
	return &v1.ListUserReply{}, nil
}
