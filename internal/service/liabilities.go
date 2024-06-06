package service

import (
	"context"
	"errors"

	v1 "core/api/v1"
	"core/internal/biz"
)

type LiabilitiesService struct {
	v1.UnimplementedLiabilitiesServer

	action *biz.LiabilityAction
}

func NewLiabilitiesService(action *biz.LiabilityAction) *LiabilitiesService {
	return &LiabilitiesService{action: action}
}

func (s *LiabilitiesService) GetLiability(ctx context.Context, req *v1.GetLiabilityRequest) (*v1.GetLiabilityReply, error) {
	if req.Id == "" {
		return &v1.GetLiabilityReply{Ok: false, Id: ""}, errors.New("id not supplied")
	}

	var lia *biz.Liability
	lia, err := s.action.FindLiabilityByID(ctx, req.Id)

	return &v1.GetLiabilityReply{Id: lia.ID, Ok: err == nil, Liability: biz.LiabilityToProtoData(lia)}, nil
}

func (s *LiabilitiesService) GetLiabilities(ctx context.Context, req *v1.GetLiabilitiesRequest) (*v1.GetLiabilitiesReply, error) {
	if req.Owner == "" {
		return &v1.GetLiabilitiesReply{Ok: false}, errors.New("id not supplied")
	}

	var lia []*biz.Liability
	var err error
	lia, err = s.action.GetLiabilities(ctx, req)

	var credit []*v1.CreditLiability
	var mortgage []*v1.MortgageLiability
	var student []*v1.StudentLiability

	for _, el := range lia {
		if el.MortgageLiability.Id != "" {
			mortgage = append(mortgage, el.MortgageLiability)
		}
		if el.CreditLiability.Id != "" {
			credit = append(credit, biz.CreditLiabilityToProto(el.CreditLiability))
		}
		if el.StudentLiability.Id != "" {
			student = append(student, el.StudentLiability)
		}
	}

	return &v1.GetLiabilitiesReply{Ok: err == nil, Credit: credit, Mortgage: mortgage, Student: student}, err
}
