package leads_service

import (
	"context"
	"errors"

	v1 "core/api/v1"
	leadsV1 "core/api/v1/b2b"
	leads_biz "core/pkg/leads/biz"
)

type CompanyService struct {
	leadsV1.UnimplementedCompaniesServer

	action *leads_biz.CompanyAction
}

func NewCompanyService(action *leads_biz.CompanyAction) *CompanyService {
	return &CompanyService{action: action}
}

func (s *CompanyService) GetCompany(ctx context.Context, req *leadsV1.GetCompanyRequest) (*leadsV1.GetCompanyReply, error) {
	if req.Id == "" {
		return &leadsV1.GetCompanyReply{Ok: false, Id: ""}, errors.New("id not supplied")
	}

	company, err := s.action.GetCompany(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &leadsV1.GetCompanyReply{
		Ok:      err == nil,
		Id:      company.ID,
		Company: leads_biz.CompanyToProtoData(company),
	}, nil
}

func (s *CompanyService) GetUSASpending(ctx context.Context, req *v1.GetUSASpendingRequest) (*v1.GetUSASpendingReply, error) {
	// TODO: Implement the logic for retrieving USA spending data
	return &v1.GetUSASpendingReply{}, nil
}

func (s *CompanyService) GetSenateLobbying(ctx context.Context, req *v1.GetSenateLobbyingRequest) (*v1.GetSenateLobbyingReply, error) {
	// TODO: Implement the logic for retrieving Senate lobbying data
	return &v1.GetSenateLobbyingReply{}, nil
}
