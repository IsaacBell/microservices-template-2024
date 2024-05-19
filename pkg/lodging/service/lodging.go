package lodging_service

import (
	"context"

	lodgingV1 "microservices-template-2024/api/v1/lodging"
	"microservices-template-2024/internal/biz"
	lodging_biz "microservices-template-2024/pkg/lodging/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyService struct {
	lodgingV1.UnimplementedLodgingServer

	propertyAction *lodging_biz.PropertyAction
	userAction     *biz.UserAction
}

func NewPropertyService(propertyAction *lodging_biz.PropertyAction, userAction *lodging_biz.UserAction) *PropertyService {
	return &PropertyService{propertyAction: propertyAction, userAction: userAction}
}

func (s *PropertyService) CreateLodging(ctx context.Context, req *lodgingV1.CreateLodgingRequest) (*lodgingV1.CreateLodgingReply, error) {
	property := lodging_biz.ProtoToPropertyData(req.Property)

	// Check if the user exists, if not, create a new user
	user, err := s.userAction.FindUserById(ctx, property.UserID)
	if err != nil {
		// Assume we missed a sync with Firebase
		user = biz.ProtoToUserData(req.Property.User)
		user, err = s.userAction.CreateUser(ctx, user)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	property.UserID = user.ID

	createdProperty, err := s.propertyAction.CreateProperty(ctx, property)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &lodgingV1.CreateLodgingReply{
		Property: lodging_biz.PropertyToProtoData(createdProperty),
	}, nil
}

func (s *PropertyService) UpdateLodging(ctx context.Context, req *lodgingV1.UpdateLodgingRequest) (*lodgingV1.UpdateLodgingReply, error) {
	property := lodging_biz.ProtoToPropertyData(req.Property)
	updatedProperty, err := s.propertyAction.UpdateProperty(ctx, property)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &lodgingV1.UpdateLodgingReply{
		Property: lodging_biz.PropertyToProtoData(updatedProperty),
	}, nil
}

func (s *PropertyService) DeleteLodging(ctx context.Context, req *lodgingV1.DeleteLodgingRequest) (*lodgingV1.DeleteLodgingReply, error) {
	err := s.propertyAction.DeleteProperty(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &lodgingV1.DeleteLodgingReply{
		Success: true,
	}, nil
}

func (s *PropertyService) GetLodging(ctx context.Context, req *lodgingV1.GetLodgingRequest) (*lodgingV1.GetLodgingReply, error) {
	property, err := s.propertyAction.GetProperty(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if property.Deleted {
		return nil, status.Error(codes.NotFound, "Property not found")
	}
	return &lodgingV1.GetLodgingReply{
		Property: lodging_biz.PropertyToProtoData(property),
	}, nil
}

func (s *PropertyService) ListLodging(ctx context.Context, req *lodgingV1.ListLodgingRequest) (*lodgingV1.ListLodgingReply, error) {
	filters := &lodging_biz.PropertyFilter{
		UserID:    req.UserId,
		Lat:       req.Lat,
		Lng:       req.Lng,
		Distance:  req.Distance,
		Area:      req.Area,
		Rooms:     req.Rooms,
		PriceGte:  req.PriceGte,
		PriceLte:  req.PriceLte,
		Sold:      req.Sold,
		Page:      int(req.Page),
		PerPage:   int(req.PerPage),
	}
	properties, err := s.propertyAction.ListProperties(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var protoProperties []*lodgingV1.Property
	for _, property := range properties {
		protoProperties = append(protoProperties, lodging_biz.PropertyToProtoData(property))
	}
	return &lodgingV1.ListLodgingReply{
		Properties: protoProperties,
	}, nil
}

func (s *PropertyService) SearchLodging(ctx context.Context, req *lodgingV1.SearchLodgingRequest) (*lodgingV1.SearchLodgingReply, error) {
	filters := &lodging_biz.PropertyFilter{
		Lat:      req.Lat,
		Lng:      req.Lng,
		Distance: req.Distance,
		Area:     req.Area,
		Rooms:    req.Rooms,
		Price:    req.Price,
		Sold:     req.Sold,
		Location: req.Location,
		Page:     int(req.Page),
		PerPage:  50,
	}
	properties, err := s.propertyAction.(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var protoProperties []*lodgingV1.Property
	for _, property := range properties {
		protoProperties = append(protoProperties, lodging_biz.PropertyToProtoData(property))
	}
	return &lodgingV1.SearchLodgingReply{
		Properties: protoProperties,
	}, nil
}

func (s *PropertyService) RealtorStats(ctx context.Context, req *lodgingV1.RealtorStatsRequest) (*lodgingV1.RealtorStatsReply, error) {
	user, err := s.userAction.FindUserById(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	stats := user.RealtorStats()
	return &lodgingV1.RealtorStatsReply{
		Stats: stats,
	}, nil
}

type UserAction struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserAction(repo UserRepo, logger log.Logger) *UserAction {
	return &UserAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserAction) FindUserById(ctx context.Context, id string) (*biz.User, error) {
	uc.log.WithContext(ctx).Infof("FindUserById: %s", id)
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserAction) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %s", user.Email)
	createdUser, err := uc.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

type UserRepo interface {
	GetByID(context.Context, string) (*biz.User, error)
	Save(context.Context, *biz.User) (*biz.User, error)
}