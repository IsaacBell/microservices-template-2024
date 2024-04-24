package service

import (
	"context"
	"io"

	pb "microservices-template-2024/api/v1/server"
	"microservices-template-2024/internal/stream"
)

type LogService struct {
	pb.UnimplementedLogServer
}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) ProduceLog(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceResponse, error) {
	if req.Message == "" {
		stream.ProduceKafkaMessage(req.Topic, string(req.Record.Value))
	} else {
		stream.ProduceKafkaMessage(req.Topic, req.Message)
	}
	return &pb.ProduceResponse{}, nil
}
func (s *LogService) ConsumeLog(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeResponse, error) {
	return &pb.ConsumeResponse{}, nil
}
func (s *LogService) ConsumeStream(req *pb.ConsumeRequest, conn pb.Log_ConsumeStreamServer) error {
	for {
		err := conn.Send(&pb.ConsumeResponse{})
		if err != nil {
			return err
		}
	}
}
func (s *LogService) ProduceStream(conn pb.Log_ProduceStreamServer) error {
	for {
		_, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = conn.Send(&pb.ProduceResponse{})
		if err != nil {
			return err
		}
	}
}
