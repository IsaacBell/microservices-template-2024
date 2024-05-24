package consultants_service

import (
	"context"

	consultantV1 "microservices-template-2024/api/v1/consultant"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
)

type ConsultantService struct {
	consultantV1.UnimplementedConsultantsServer

	action *consultants_biz.ConsultantAction
}

func NewConsultantService(action *consultants_biz.ConsultantAction) *ConsultantService {
	return &ConsultantService{action: action}
}

func (s *ConsultantService) GetConsultant(ctx context.Context, req *consultantV1.GetConsultantRequest) (*consultantV1.GetConsultantReply, error) {
	consultant, err := s.action.GetConsultant(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &consultantV1.GetConsultantReply{
		Ok:         err == nil,
		Consultant: consultants_biz.ConsultantToProtoData(consultant),
	}, nil
}

func (s *ConsultantService) CreateConsultant(ctx context.Context, req *consultantV1.CreateConsultantRequest) (*consultantV1.CreateConsultantReply, error) {
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	createdConsultant, err := s.action.CreateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}
	return &consultantV1.CreateConsultantReply{
		Ok:         err == nil,
		Consultant: consultants_biz.ConsultantToProtoData(createdConsultant),
	}, nil
}

func (s *ConsultantService) UpdateConsultant(ctx context.Context, req *consultantV1.UpdateConsultantRequest) (*consultantV1.UpdateConsultantReply, error) {
	consultant := consultants_biz.ProtoToConsultantData(req.Consultant)
	updatedConsultant, err := s.action.UpdateConsultant(ctx, consultant)
	if err != nil {
		return nil, err
	}
	return &consultantV1.UpdateConsultantReply{
		Ok:         err == nil,
		Consultant: consultants_biz.ConsultantToProtoData(updatedConsultant),
	}, nil
}

func (s *ConsultantService) DeleteConsultant(ctx context.Context, req *consultantV1.DeleteConsultantRequest) (*consultantV1.DeleteConsultantReply, error) {
	err := s.action.DeleteConsultant(ctx, req.Id)
	return &consultantV1.DeleteConsultantReply{
		Ok: err == nil,
	}, err
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
	consultants, err := s.action.ListConsultants(ctx, filters)
	if err != nil {
		return nil, err
	}
	protoConsultants := make([]*consultantV1.Consultant, len(consultants))
	for i, consultant := range consultants {
		protoConsultants[i] = consultants_biz.ConsultantToProtoData(consultant)
	}
	return &consultantV1.ListConsultantsReply{
		Ok:          err == nil,
		Consultants: protoConsultants,
	}, nil
}
