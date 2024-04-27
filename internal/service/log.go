package service

import (
	"context"
	"io"

	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/stream"
)

type LogService struct {
	v1.UnimplementedLogServer
}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) ProduceLog(ctx context.Context, req *v1.ProduceRequest) (*v1.ProduceResponse, error) {
	if req.Message == "" {
		stream.ProduceKafkaMessage(req.Topic, string(req.Record.Value))
	} else {
		stream.ProduceKafkaMessage(req.Topic, req.Message)
	}
	return &v1.ProduceResponse{}, nil
}
func (s *LogService) ConsumeLog(ctx context.Context, req *v1.ConsumeRequest) (*v1.ConsumeResponse, error) {
	return &v1.ConsumeResponse{}, nil
}
func (s *LogService) ConsumeStream(req *v1.ConsumeRequest, conn v1.Log_ConsumeStreamServer) error {
	for {
		err := conn.Send(&v1.ConsumeResponse{})
		if err != nil {
			return err
		}
	}
}
func (s *LogService) ProduceStream(conn v1.Log_ProduceStreamServer) error {
	for {
		rec, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// conn.Send undefined (type "microservices-template-2024/api/v1".Log_ProduceStreamServer has no field or method Send)compilerMissingFieldOrMethod
		err = conn.SendMsg(&v1.ProduceResponse{Offset: rec.Record.Offset})
		if err != nil {
			return err
		}
	}
}
