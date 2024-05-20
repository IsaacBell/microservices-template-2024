package lodging_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	lodgingV1 "microservices-template-2024/api/v1/lodging"
	lodging_biz "microservices-template-2024/pkg/lodging/biz"

	cache "microservices-template-2024/pkg/cache"
	stream "microservices-template-2024/pkg/stream"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyService struct {
	lodgingV1.UnimplementedLodgingServer

	propertyAction *lodging_biz.PropertyAction
	userAction     *lodging_biz.UserAction
}

func isPresent(val interface{}) bool {
	return val != nil && val != 0 && val != ""
}

func filtersFromReq(req interface{}) map[string]interface{} {
	filters := make(map[string]interface{})

	switch r := req.(type) {
	case *lodgingV1.ListLodgingRequest:
		if isPresent(r.UserId) {
			filters["userId"] = r.UserId
		}
		if isPresent(r.Lat) {
			filters["lat"] = r.Lat
		}
		if isPresent(r.Lng) {
			filters["lng"] = r.Lng
		}
		if isPresent(r.Distance) {
			filters["distance"] = r.Distance
		}
		if isPresent(r.Area) {
			filters["area"] = r.Area
		}
		if isPresent(r.Rooms) {
			filters["rooms"] = r.Rooms
		}
		if isPresent(r.PriceGte) {
			filters["priceGte"] = r.PriceGte
		}
		if isPresent(r.PriceLte) {
			filters["priceLte"] = r.PriceLte
		}
		if isPresent(r.Sold) {
			filters["sold"] = r.Sold
		}
		if isPresent(r.Page) {
			filters["page"] = r.Page
		}
		if isPresent(r.PerPage) {
			filters["perPage"] = r.PerPage
		}
	case *lodgingV1.SearchLodgingRequest:
		if isPresent(r.Lat) {
			filters["lat"] = r.Lat
		}
		if isPresent(r.Lng) {
			filters["lng"] = r.Lng
		}
		if isPresent(r.Distance) {
			filters["distance"] = r.Distance
		}
		if isPresent(r.Area) {
			filters["area"] = r.Area
		}
		if isPresent(r.Rooms) {
			filters["rooms"] = r.Rooms
		}
		if isPresent(r.Price) {
			filters["price"] = r.Price
		}
		if isPresent(r.Sold) {
			filters["sold"] = r.Sold
		}
		if isPresent(r.Location) {
			filters["location"] = r.Location
		}
		if isPresent(r.Page) {
			filters["page"] = r.Page
		}
	}

	return filters
}

func NewPropertyService(propertyAction *lodging_biz.PropertyAction) *PropertyService {
	return &PropertyService{propertyAction: propertyAction}
}

func (s *PropertyService) CreateLodging(ctx context.Context, req *lodgingV1.CreateLodgingRequest) (*lodgingV1.CreateLodgingReply, error) {
	property := lodging_biz.ProtoToPropertyData(req.Property)

	// Check if the user exists, if not, create a new user
	user, err := s.userAction.FindUserById(ctx, property.UserID)
	if err != nil || user == nil {
		// Assume we missed a sync with Firebase
		user = lodging_biz.ProtoToUserData(req.Property.User)
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

	cache.Cache(ctx).Set("property:"+createdProperty.ID, createdProperty, time.Hour*24)
	stream.ProduceKafkaMessage("channel/lodging", "New Property: "+createdProperty.ID)

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

	cache.Cache(ctx).Set("property:"+updatedProperty.ID, updatedProperty, time.Hour*24)
	stream.ProduceKafkaMessage("channel/lodging", "Update Property: "+updatedProperty.ID + ":" + updatedProperty.Address)

	return &lodgingV1.UpdateLodgingReply{
		Property: lodging_biz.PropertyToProtoData(updatedProperty),
	}, nil
}

func (s *PropertyService) DeleteLodging(ctx context.Context, req *lodgingV1.DeleteLodgingRequest) (*lodgingV1.DeleteLodgingReply, error) {
	err := s.propertyAction.DeleteProperty(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func() {
		err := cache.Cache(ctx).Del(req.Id).Err()
		if err != nil {
			fmt.Printf("Failed to delete cache entry for property %d: %v \n", req.Id, err)
		}
	}()
	stream.ProduceKafkaMessage("channel/lodging", "Deleted Property: "+req.Id)

	return &lodgingV1.DeleteLodgingReply{
		Success: true,
	}, nil
}

func (s *PropertyService) GetLodging(ctx context.Context, req *lodgingV1.GetLodgingRequest) (*lodgingV1.GetLodgingReply, error) {
	cacheKey := "property:"+req.Id
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		// Cache hit, deserialize the cached data
		var property lodgingV1.Property
		if err := json.Unmarshal([]byte(cached), &property); err == nil {
			return &lodgingV1.GetLodgingReply{
				Property: &property,
			}, nil
		}
	}

	property, err := s.propertyAction.GetProperty(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if property.Deleted {
		return nil, status.Error(codes.NotFound, "Property not found")
	}
	stream.ProduceKafkaMessage("channel/lodging", "Cache Property: "+req.Id)

	cache.Cache(ctx).Set(cacheKey, property, time.Hour*24)

	return &lodgingV1.GetLodgingReply{
		Property: lodging_biz.PropertyToProtoData(property),
	}, nil
}

func (s *PropertyService) ListLodging(ctx context.Context, req *lodgingV1.ListLodgingRequest) (*lodgingV1.ListLodgingReply, error) {
	filters := filtersFromReq(req)
	properties, err := s.propertyAction.ListLodging(ctx, filters)
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
	filters := filtersFromReq(req)
	cacheKey := fmt.Sprintf("lodging_search:%v", filters)

	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		// Cache hit, deserialize the cached data
		var protoProperties []*lodgingV1.Property
		if err := json.Unmarshal([]byte(cached), &protoProperties); err == nil {
			return &lodgingV1.SearchLodgingReply{
				Properties: protoProperties,
			}, nil
		}
	}

	properties, err := s.propertyAction.SearchProperties(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoProperties := make([]*lodgingV1.Property, len(properties))
	for _, property := range properties {
		protoProperties = append(protoProperties, lodging_biz.PropertyToProtoData(property))
	}

	if data, err := json.Marshal(protoProperties); err == nil {
		cache.Cache(ctx).Set(cacheKey, string(data), time.Hour*12)
	}

	stream.ProduceKafkaMessage("channel/lodging", "Cached: "+cacheKey)

	return &lodgingV1.SearchLodgingReply{
		Properties: protoProperties,
	}, nil
}

func (s *PropertyService) RealtorStats(ctx context.Context, req *lodgingV1.RealtorStatsRequest) (*lodgingV1.RealtorStatsReply, error) {
	cacheKey := "realtor_stats:"+req.UserId
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		// Cache hit, deserialize the cached data
		var stats map[string]int64
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &lodgingV1.RealtorStatsReply{
				Stats: stats,
			}, nil
		}
	}

	stats, err := s.propertyAction.GetRealtorStats(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	cache.Cache(ctx).Set(cacheKey, stats, time.Hour*12)
	stream.ProduceKafkaMessage("channel/lodging", "Cached: "+cacheKey)

	return &lodgingV1.RealtorStatsReply{
		Stats: stats,
	}, nil
}
