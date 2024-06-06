package service

import (
	"context"
	"fmt"
	"io"

	v1 "core/api/v1"
	"core/pkg/stream"
)

type LogService struct {
	v1.UnimplementedLogServer
}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) Produce(ctx context.Context, req *v1.ProduceRequest) (*v1.ProduceResponse, error) {
	if req.Message == "" {
		stream.ProduceKafkaMessage(req.Topic, string(req.Record.Value))
	} else {
		stream.ProduceKafkaMessage(req.Topic, req.Message)
	}
	return &v1.ProduceResponse{Ok: true}, nil
}

func (s *LogService) Consume(ctx context.Context, req *v1.ConsumeRequest) (*v1.ConsumeResponse, error) {
	return &v1.ConsumeResponse{}, nil
}

func (s *LogService) ConsumeStream(req *v1.ConsumeRequest, conn v1.Log_ConsumeStreamServer) error {
	ctx, cancel := stream.StartKafkaConsumer(req.Topic, "core", func(msg string) {
		fmt.Println("Received Kafka Msg: [", req.Topic, "] ", msg)
		err := conn.Send(&v1.ConsumeResponse{
			Record: &v1.Record{
				Value:  []byte(msg),
				Offset: 0,
			},
		})
		if err != nil {
			return
		}
	})
	defer cancel()

	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				cancel()
				return
			default:
				// Do nothing
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
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

		// conn.Send undefined (type "core/api/v1".Log_ProduceStreamServer has no field or method Send)compilerMissingFieldOrMethod
		err = conn.SendMsg(&v1.ProduceResponse{Offset: rec.Record.Offset})
		if err != nil {
			return err
		}
	}
}
