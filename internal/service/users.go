package service

import (
	"context"

	pb "microservices-template-2024/api/v1/server"
	// log "microservices-template-2024/internal/service/log"
)

type UsersService struct {
	pb.UnimplementedUsersServer
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (s *UsersService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	
	return &pb.CreateUserReply{}, nil
}
func (s *UsersService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return &pb.UpdateUserReply{}, nil
}
func (s *UsersService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	return &pb.DeleteUserReply{}, nil
}
func (s *UsersService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	u := pb.User{Email: *req.Email, Id: *req.Id}
	return &pb.GetUserReply{}, nil
}
func (s *UsersService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}
