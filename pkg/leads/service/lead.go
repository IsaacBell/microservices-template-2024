package leads_service

import (
	"context"
	"errors"

	v1 "microservices-template-2024/api/v1"
	leadsV1 "microservices-template-2024/api/v1/b2b"
	leads_biz "microservices-template-2024/pkg/leads/biz"
)

type LeadService struct {
	leadsV1.UnimplementedLeadsServer

	action *leads_biz.LeadAction
}

func NewLeadService(action *leads_biz.LeadAction) *LeadService {
	return &LeadService{action: action}
}

func (s *LeadService) GetLead(ctx context.Context, req *leadsV1.GetLeadRequest) (*leadsV1.GetLeadReply, error) {
	if req.Id == "" {
		return &leadsV1.GetLeadReply{Ok: false, Id: ""}, errors.New("id not supplied")
	}

	lead, err := s.action.GetLead(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &leadsV1.GetLeadReply{
		Ok:   err == nil,
		Id:   lead.ID,
		Lead: leads_biz.LeadToProtoData(lead),
	}, nil
}

func (s *LeadService) GetUSASpending(ctx context.Context, req *v1.GetUSASpendingRequest) (*v1.GetUSASpendingReply, error) {
	// TODO: Implement the logic for retrieving USA spending data
	return &v1.GetUSASpendingReply{}, nil
}

func (s *LeadService) GetSenateLobbying(ctx context.Context, req *v1.GetSenateLobbyingRequest) (*v1.GetSenateLobbyingReply, error) {
	// TODO: Implement the logic for retrieving Senate lobbying data
	return &v1.GetSenateLobbyingReply{}, nil
}
