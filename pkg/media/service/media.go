package service

import (
	"context"

	pb "core/api/v1"
)

type MediaService struct {
	pb.UnimplementedMediaServer
}

func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) CreateMedia(ctx context.Context, req *pb.CreateMediaRequest) (*pb.CreateMediaReply, error) {
	return &pb.CreateMediaReply{}, nil
}
func (s *MediaService) UpdateMedia(ctx context.Context, req *pb.UpdateMediaRequest) (*pb.UpdateMediaReply, error) {
	return &pb.UpdateMediaReply{}, nil
}
func (s *MediaService) DeleteMedia(ctx context.Context, req *pb.DeleteMediaRequest) (*pb.DeleteMediaReply, error) {
	return &pb.DeleteMediaReply{}, nil
}
func (s *MediaService) GetMedia(ctx context.Context, req *pb.GetMediaRequest) (*pb.GetMediaReply, error) {
	return &pb.GetMediaReply{}, nil
}
func (s *MediaService) ListMedia(ctx context.Context, req *pb.ListMediaRequest) (*pb.ListMediaReply, error) {
	return &pb.ListMediaReply{}, nil
}
