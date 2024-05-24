package consultants_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	consultantV1 "microservices-template-2024/api/v1/consultants"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"

	cache "microservices-template-2024/pkg/cache"
	stream "microservices-template-2024/pkg/stream"
)

type ConsultantService struct {
	consultantV1.UnimplementedConsultantsServer

	action *consultants_biz.ConsultantAction
}

func NewConsultantService(action *consultants_biz.ConsultantAction) *ConsultantService {
	return &ConsultantService{action: action}
}

func (s *ConsultantService) GetConsultant(ctx context.Context, req *consultantV1.GetConsultantRequest) (*consultantV1.GetConsultantReply, error) {
	cacheKey := "consultant:" + req.Id
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		var consultant consultantV1.Consultant
		if err := json.Unmarshal([]byte(cached), &consultant); err == nil {
			return &consultantV1.GetConsultantReply{
				Ok:         true,
				Consultant: &consultant,
			}, nil
		}
	}

	consultant, err := s.action.GetConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	go func() {
		data, _ := json.Marshal(consultants_biz.ConsultantToProtoData(consultant))
		cache.Cache(ctx).Set(cacheKey, string(data), time.Hour*24)
		stream.ProduceKafkaMessage("channel/consultants", "Cached consultant: "+req.Id)
	}()

	return &consultantV1.GetConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(consultant),
	}, nil
}

func (s *ConsultantService) CreateConsultant(ctx context.Context, req *consultantV1.CreateConsultantRequest) (*consultantV1.CreateConsultantReply, error) {
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	createdConsultant, err := s.action.CreateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}

	go func() {
		data, _ := json.Marshal(consultants_biz.ConsultantToProtoData(createdConsultant))
		cache.Cache(ctx).Set("consultant:"+createdConsultant.ID, string(data), time.Hour*24)
		stream.ProduceKafkaMessage("channel/consultants", "New consultant: "+createdConsultant.ID)
	}()

	return &consultantV1.CreateConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(createdConsultant),
	}, nil
}

func (s *ConsultantService) UpdateConsultant(ctx context.Context, req *consultantV1.UpdateConsultantRequest) (*consultantV1.UpdateConsultantReply, error) {
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	updatedConsultant, err := s.action.UpdateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}

	go func() {
		data, _ := json.Marshal(consultants_biz.ConsultantToProtoData(updatedConsultant))
		cache.Cache(ctx).Set("consultant:"+updatedConsultant.ID, string(data), time.Hour*24)
		stream.ProduceKafkaMessage("channel/consultants", "Updated consultant: "+updatedConsultant.ID)
	}()

	return &consultantV1.UpdateConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(updatedConsultant),
	}, nil
}

func (s *ConsultantService) DeleteConsultant(ctx context.Context, req *consultantV1.DeleteConsultantRequest) (*consultantV1.DeleteConsultantReply, error) {
	err := s.action.DeleteConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cache.Cache(ctx).Del("consultant:" + req.Id).Err()
		if err != nil {
			fmt.Printf("Failed to delete cache entry for consultant %s: %v\n", req.Id, err)
		}
		stream.ProduceKafkaMessage("channel/consultants", "Deleted consultant: "+req.Id)
	}()

	return &consultantV1.DeleteConsultantReply{
		Ok: true,
	}, nil
}

func (s *ConsultantService) ListConsultants(ctx context.Context, req *consultantV1.ListConsultantsRequest) (*consultantV1.ListConsultantsReply, error) {
	filters := make(map[string]interface{})
	if req.UserId != "" {
		filters["user_id"] = req.UserId
	}
	if len(req.Specializations) > 0 {
		filters["specializations"] = req.Specializations
	}
	if len(req.Languages) > 0 {
		filters["languages"] = req.Languages
	}

	cacheKey := fmt.Sprintf("consultants_list:%v", filters)
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		var protoConsultants []*consultantV1.Consultant
		if err := json.Unmarshal([]byte(cached), &protoConsultants); err == nil {
			return &consultantV1.ListConsultantsReply{
				Ok:          true,
				Consultants: protoConsultants,
			}, nil
		}
	}

	consultants, err := s.action.ListConsultants(ctx, filters)
	if err != nil {
		return nil, err
	}

	protoConsultants := make([]*consultantV1.Consultant, len(consultants))
	for i, consultant := range consultants {
		protoConsultants[i] = consultants_biz.ConsultantToProtoData(consultant)
	}

	go func() {
		data, _ := json.Marshal(protoConsultants)
		cache.Cache(ctx).Set(cacheKey, string(data), time.Hour*12)
		stream.ProduceKafkaMessage("channel/consultants", "Cached consultants list: "+cacheKey)
	}()

	return &consultantV1.ListConsultantsReply{
		Ok:          true,
		Consultants: protoConsultants,
	}, nil
}
