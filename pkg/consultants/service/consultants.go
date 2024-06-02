package consultants_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	communicationsV1 "microservices-template-2024/api/v1/communications"
	consultantsV1 "microservices-template-2024/api/v1/consultants"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"

	cache "microservices-template-2024/pkg/cache"
)

type ConsultantService struct {
	consultantsV1.UnimplementedConsultantsServer

	action *consultants_biz.ConsultantAction
}

type Action interface {
	CreateConsultant(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error)
	DeleteConsultant(ctx context.Context, id string) error
	GetConsultant(ctx context.Context, id string) (*consultants_biz.Consultant, error)
	ListConsultants(ctx context.Context, filters map[string]interface{}) ([]*consultants_biz.Consultant, error)
	SendComm(ctx context.Context, c *consultants_biz.Communication) (*consultants_biz.Communication, error)
	UpdateConsultant(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error)
}

func NewConsultantService(action *consultants_biz.ConsultantAction) *ConsultantService {
	return &ConsultantService{action: action}
}

func (s *ConsultantService) GetConsultant(ctx context.Context, req *consultantsV1.GetConsultantRequest) (*consultantsV1.GetConsultantReply, error) {
	cacheKey := "consultant:" + req.Id
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		var consultant consultantsV1.Consultant
		if err := json.Unmarshal([]byte(cached), &consultant); err == nil {
			return &consultantsV1.GetConsultantReply{
				Ok:         true,
				Consultant: &consultant,
			}, nil
		}
	}

	consultant, err := s.action.GetConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	cache.CacheRecord("consultant", cacheKey, req.Id, consultant)

	return &consultantsV1.GetConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(consultant),
	}, nil
}

func (s *ConsultantService) CreateConsultant(ctx context.Context, req *consultantsV1.CreateConsultantRequest) (*consultantsV1.CreateConsultantReply, error) {
	if req.Consultant == nil {
		return &consultantsV1.CreateConsultantReply{Consultant: nil}, errors.New("id not supplied")
	}

	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	createdConsultant, err := s.action.CreateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}

	cache.CacheRecord("consultant", "consultant:"+createdConsultant.ID, createdConsultant.ID, createdConsultant)

	return &consultantsV1.CreateConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(createdConsultant),
	}, nil
}

func (s *ConsultantService) UpdateConsultant(ctx context.Context, req *consultantsV1.UpdateConsultantRequest) (*consultantsV1.UpdateConsultantReply, error) {
	if req.Consultant == nil {
		return &consultantsV1.UpdateConsultantReply{Consultant: nil}, errors.New("id not supplied")
	}
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	updatedConsultant, err := s.action.UpdateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}

	cache.CacheRecord("consultant", "consultant:"+updatedConsultant.ID, updatedConsultant.ID, updatedConsultant)

	return &consultantsV1.UpdateConsultantReply{
		Ok:         true,
		Consultant: consultants_biz.ConsultantToProtoData(updatedConsultant),
	}, nil
}

func (s *ConsultantService) DeleteConsultant(ctx context.Context, req *consultantsV1.DeleteConsultantRequest) (*consultantsV1.DeleteConsultantReply, error) {
	if req.Id == "" {
		return &consultantsV1.DeleteConsultantReply{Ok: false}, errors.New("id not supplied")
	}
	err := s.action.DeleteConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	cacheKey := "consultant:" + req.Id
	cache.Delete(ctx, cacheKey, "consultant")

	return &consultantsV1.DeleteConsultantReply{
		Ok: true,
	}, nil
}

func (s *ConsultantService) ListConsultants(ctx context.Context, req *consultantsV1.ListConsultantsRequest) (*consultantsV1.ListConsultantsReply, error) {
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
		var protoConsultants []*consultantsV1.Consultant
		if err := json.Unmarshal([]byte(cached), &protoConsultants); err == nil {
			return &consultantsV1.ListConsultantsReply{
				Ok:          true,
				Consultants: protoConsultants,
			}, nil
		}
	}

	consultants, err := s.action.ListConsultants(ctx, filters)
	if err != nil {
		return nil, err
	}

	protoConsultants := make([]*consultantsV1.Consultant, len(consultants))
	for i, consultant := range consultants {
		protoConsultants[i] = consultants_biz.ConsultantToProtoData(consultant)
	}

	cache.CacheRecord("consultants_list", cacheKey, cacheKey, protoConsultants)

	return &consultantsV1.ListConsultantsReply{
		Ok:          true,
		Consultants: protoConsultants,
	}, nil
}

func (s *ConsultantService) SendComm(
	ctx context.Context,
	req *consultantsV1.SendCommsRequest,
) (*consultantsV1.SendCommsReply, error) {
	if req.Comm.UserId == "" {
		return &consultantsV1.SendCommsReply{Ok: false, Ack: nil}, errors.New("user id not supplied")
	}

	res, err := s.action.SendComm(ctx, consultants_biz.ProtoToCommunicationData(req.Comm))
	if err != nil {
		return nil, err
	}
	fmt.Println("cid: ", req.Comm.UserId, ", rid: ", req.Comm.RecipientId)
	return &consultantsV1.SendCommsReply{
		Ok: err == nil,
		Ack: &communicationsV1.Ack{
			Id:              strconv.FormatUint(uint64(0), 10),
			From:            req.Comm.UserId,
			UserId:          req.Comm.UserId,
			Msg:             req.Comm.Msg,
			RecipientsCount: 1,
			Recipients:      []string{res.RecipientID},
		},
	}, nil
}
