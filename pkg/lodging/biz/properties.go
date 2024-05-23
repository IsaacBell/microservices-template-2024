package lodging_biz

import (
	"context"
	"fmt"
	lodgingV1 "microservices-template-2024/api/v1/lodging"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	ID             string                 `gorm:"primaryKey" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title          string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description    string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Address        string                 `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Area           int32                  `protobuf:"varint,5,opt,name=area,proto3" json:"area,omitempty"`
	Rooms          int32                  `protobuf:"varint,6,opt,name=rooms,proto3" json:"rooms,omitempty"`
	Price          int32                  `protobuf:"varint,7,opt,name=price,proto3" json:"price,omitempty"`
	Sold           bool                   `protobuf:"varint,8,opt,name=sold,proto3" json:"sold,omitempty"`
	Deleted        bool                   `protobuf:"varint,9,opt,name=deleted,proto3" json:"deleted,omitempty"`
	Images         []string               `protobuf:"bytes,10,rep,name=images,proto3" json:"images,omitempty"`
	PreviewImages  []string               `protobuf:"bytes,11,rep,name=preview_images,json=previewImages,proto3" json:"preview_images,omitempty"`
	Location       *Location              `protobuf:"bytes,13,opt,name=location,proto3" json:"location,omitempty"`
	UserID         string                 `protobuf:"bytes,14,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,16,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,17,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Vendor         *VendorType            `protobuf:"bytes,18,opt,name=vendor,proto3" json:"vendor,omitempty"`
	Equipments     []*EquipmentType       `protobuf:"bytes,19,rep,name=equipments,proto3" json:"equipments,omitempty"`
	Specifications []*SpecificationType   `protobuf:"bytes,20,rep,name=specifications,proto3" json:"specifications,omitempty"`
	Reviews        []*ReviewType          `protobuf:"bytes,21,rep,name=reviews,proto3" json:"reviews,omitempty"`
	ReviewStats    *ReviewStatsType       `protobuf:"bytes,22,opt,name=review_stats,json=reviewStats,proto3" json:"review_stats,omitempty"`
}

type PropertyIndex struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Address        string                 `json:"address"`
	Area           int32                  `json:"area"`
	Rooms          int32                  `json:"rooms"`
	Price          int32                  `json:"price"`
	Sold           bool                   `json:"sold"`
	Deleted        bool                   `json:"deleted"`
	Images         []string               `json:"images"`
	PreviewImages  []string               `json:"preview_images"`
	Location       *Location              `json:"location"`
	UserID         string                 `json:"user_id"`
	CreatedAt      *timestamppb.Timestamp `json:"created_at"`
	UpdatedAt      *timestamppb.Timestamp `json:"updated_at"`
	Vendor         *VendorType            `json:"vendor"`
	Equipments     []*EquipmentType       `json:"equipments"`
	Specifications []*SpecificationType   `json:"specifications"`
	Reviews        []*ReviewType          `json:"reviews"`
	ReviewStats    *ReviewStatsType       `json:"review_stats"`
}

func LocationToProtoData(location *Location) *lodgingV1.Location {
	if location == nil {
		return nil
	}

	return &lodgingV1.Location{
		Type:        location.Type,
		Coordinates: location.Coordinates,
	}
}

func ProtoToLocationData(location *lodgingV1.Location) *Location {
	if location == nil {
		return nil
	}

	return &Location{
		Type:        location.Type,
		Coordinates: location.Coordinates,
	}
}

func VendorTypeToProtoData(vendor *VendorType) *lodgingV1.VendorType {
	if vendor == nil {
		return nil
	}

	return &lodgingV1.VendorType{
		Name:          vendor.Name,
		Img:           vendor.Img,
		MemberSince:   vendor.MemberSince,
		Languages:     vendor.Languages,
		ResponseRate:  vendor.ResponseRate,
		ResponseTime:  vendor.ResponseTime,
		Location:      vendor.Location,
		BoatName:      vendor.BoatName,
		BoatGuests:    vendor.BoatGuests,
		BoatCabins:    vendor.BoatCabins,
		BoatBathrooms: vendor.BoatBathrooms,
		TotalReview:   vendor.TotalReview,
	}
}

func ProtoToVendorTypeData(vendor *lodgingV1.VendorType) *VendorType {
	if vendor == nil {
		return nil
	}

	return &VendorType{
		Name:          vendor.Name,
		Img:           vendor.Img,
		MemberSince:   vendor.MemberSince,
		Languages:     vendor.Languages,
		ResponseRate:  vendor.ResponseRate,
		ResponseTime:  vendor.ResponseTime,
		Location:      vendor.Location,
		BoatName:      vendor.BoatName,
		BoatGuests:    vendor.BoatGuests,
		BoatCabins:    vendor.BoatCabins,
		BoatBathrooms: vendor.BoatBathrooms,
		TotalReview:   vendor.TotalReview,
	}
}

func (p *Property) ToIndex() *PropertyIndex {
	return &PropertyIndex{
		ID:             p.ID,
		Title:          p.Title,
		Description:    p.Description,
		Address:        p.Address,
		Area:           p.Area,
		Rooms:          p.Rooms,
		Price:          p.Price,
		Sold:           p.Sold,
		Deleted:        p.Deleted,
		Images:         p.Images,
		PreviewImages:  p.PreviewImages,
		Location:       p.Location,
		UserID:         p.UserID,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Vendor:         p.Vendor,
		Equipments:     p.Equipments,
		Specifications: p.Specifications,
		Reviews:        p.Reviews,
		ReviewStats:    p.ReviewStats,
	}
}

func EquipmentTypeToProtoData(equipment *EquipmentType) *lodgingV1.EquipmentType {
	if equipment == nil {
		return nil
	}

	return &lodgingV1.EquipmentType{
		Img:  equipment.Img,
		Name: equipment.Name,
	}
}

func EquipmentTypesToProtoData(equipments []*EquipmentType) []*lodgingV1.EquipmentType {
	var protos []*lodgingV1.EquipmentType
	for _, equipment := range equipments {
		protos = append(protos, EquipmentTypeToProtoData(equipment))
	}
	return protos
}

func ProtoToEquipmentTypeData(equipment *lodgingV1.EquipmentType) *EquipmentType {
	if equipment == nil {
		return nil
	}

	return &EquipmentType{
		Img:  equipment.Img,
		Name: equipment.Name,
	}
}

func ProtoToEquipmentTypesData(equipments []*lodgingV1.EquipmentType) []*EquipmentType {
	var types []*EquipmentType
	for _, equipment := range equipments {
		types = append(types, ProtoToEquipmentTypeData(equipment))
	}
	return types
}

func SpecificationTypeToProtoData(specification *SpecificationType) *lodgingV1.SpecificationType {
	if specification == nil {
		return nil
	}

	return &lodgingV1.SpecificationType{
		Name:    specification.Name,
		Details: specification.Details,
	}
}

func SpecificationTypesToProtoData(specifications []*SpecificationType) []*lodgingV1.SpecificationType {
	var protos []*lodgingV1.SpecificationType
	for _, specification := range specifications {
		protos = append(protos, SpecificationTypeToProtoData(specification))
	}
	return protos
}

func ProtoToSpecificationTypeData(specification *lodgingV1.SpecificationType) *SpecificationType {
	if specification == nil {
		return nil
	}

	return &SpecificationType{
		Name:    specification.Name,
		Details: specification.Details,
	}
}

func ProtoToSpecificationTypesData(specifications []*lodgingV1.SpecificationType) []*SpecificationType {
	var types []*SpecificationType
	for _, specification := range specifications {
		types = append(types, ProtoToSpecificationTypeData(specification))
	}
	return types
}

func ReviewTypeToProtoData(review *ReviewType) *lodgingV1.ReviewType {
	if review == nil {
		return nil
	}

	return &lodgingV1.ReviewType{
		Avatar:   review.Avatar,
		Name:     review.Name,
		Date:     review.Date,
		Rating:   review.Rating,
		Location: review.Location,
		Review:   review.Review,
	}
}

func ReviewTypesToProtoData(reviews []*ReviewType) []*lodgingV1.ReviewType {
	var protos []*lodgingV1.ReviewType
	for _, review := range reviews {
		protos = append(protos, ReviewTypeToProtoData(review))
	}
	return protos
}

func ProtoToReviewTypeData(review *lodgingV1.ReviewType) *ReviewType {
	if review == nil {
		return nil
	}

	return &ReviewType{
		Avatar:   review.Avatar,
		Name:     review.Name,
		Date:     review.Date,
		Rating:   review.Rating,
		Location: review.Location,
		Review:   review.Review,
	}
}

func ProtoToReviewTypesData(reviews []*lodgingV1.ReviewType) []*ReviewType {
	var types []*ReviewType
	for _, review := range reviews {
		types = append(types, ProtoToReviewTypeData(review))
	}
	return types
}

func ReviewBarTypeToProtoData(reviewBar *ReviewBarType) *lodgingV1.ReviewBarType {
	if reviewBar == nil {
		return nil
	}

	return &lodgingV1.ReviewBarType{
		Count:   reviewBar.Count,
		Percent: reviewBar.Percent,
	}
}

func ReviewBarTypesToProtoData(reviewBars []*ReviewBarType) []*lodgingV1.ReviewBarType {
	var protos []*lodgingV1.ReviewBarType
	for _, reviewBar := range reviewBars {
		protos = append(protos, ReviewBarTypeToProtoData(reviewBar))
	}
	return protos
}

func ProtoToReviewBarTypeData(reviewBar *lodgingV1.ReviewBarType) *ReviewBarType {
	if reviewBar == nil {
		return nil
	}

	return &ReviewBarType{
		Count:   reviewBar.Count,
		Percent: reviewBar.Percent,
	}
}

func ProtoToReviewBarTypesData(reviewBars []*lodgingV1.ReviewBarType) []*ReviewBarType {
	var types []*ReviewBarType
	for _, reviewBar := range reviewBars {
		types = append(types, ProtoToReviewBarTypeData(reviewBar))
	}
	return types
}

func ReviewStatsTypeToProtoData(reviewStats *ReviewStatsType) *lodgingV1.ReviewStatsType {
	if reviewStats == nil {
		return nil
	}

	return &lodgingV1.ReviewStatsType{
		TotalReviews:  reviewStats.TotalReviews,
		AverageRating: reviewStats.AverageRating,
		Stars:         ReviewBarTypesToProtoData(reviewStats.Stars),
	}
}

func ProtoToReviewStatsTypeData(reviewStats *lodgingV1.ReviewStatsType) *ReviewStatsType {
	if reviewStats == nil {
		return nil
	}

	return &ReviewStatsType{
		TotalReviews:  reviewStats.TotalReviews,
		AverageRating: reviewStats.AverageRating,
		Stars:         ProtoToReviewBarTypesData(reviewStats.Stars),
	}
}

func (p *Property) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func PropertyToProtoData(property *Property) *lodgingV1.Property {
	if property == nil {
		return nil
	}

	return &lodgingV1.Property{
		Id:             property.ID,
		Title:          property.Title,
		Description:    property.Description,
		Address:        property.Address,
		Area:           property.Area,
		Rooms:          property.Rooms,
		Price:          property.Price,
		Sold:           property.Sold,
		Deleted:        property.Deleted,
		Images:         property.Images,
		PreviewImages:  property.PreviewImages,
		Location:       LocationToProtoData(property.Location),
		UserId:         property.UserID,
		CreatedAt:      property.CreatedAt,
		UpdatedAt:      property.UpdatedAt,
		Vendor:         VendorTypeToProtoData(property.Vendor),
		Equipments:     EquipmentTypesToProtoData(property.Equipments),
		Specifications: SpecificationTypesToProtoData(property.Specifications),
		Reviews:        ReviewTypesToProtoData(property.Reviews),
		ReviewStats:    ReviewStatsTypeToProtoData(property.ReviewStats),
	}
}

func ProtoToPropertyData(input *lodgingV1.Property) *Property {
	property := &Property{
		ID:             input.Id,
		Title:          input.Title,
		Description:    input.Description,
		Address:        input.Address,
		Area:           input.Area,
		Rooms:          input.Rooms,
		Price:          input.Price,
		Sold:           input.Sold,
		Deleted:        input.Deleted,
		Images:         input.Images,
		PreviewImages:  input.PreviewImages,
		Location:       ProtoToLocationData(input.Location),
		UserID:         input.UserId,
		CreatedAt:      input.CreatedAt,
		UpdatedAt:      input.UpdatedAt,
		Vendor:         ProtoToVendorTypeData(input.Vendor),
		Equipments:     ProtoToEquipmentTypesData(input.Equipments),
		Specifications: ProtoToSpecificationTypesData(input.Specifications),
		Reviews:        ProtoToReviewTypesData(input.Reviews),
		ReviewStats:    ProtoToReviewStatsTypeData(input.ReviewStats),
	}

	return property
}

type PropertyRepo interface {
	Get(context.Context, string) (*Property, error)
	Save(context.Context, *Property) (*Property, error)
	Stats(context.Context, string) (map[string]int64, error)
	Update(context.Context, *Property) (*Property, error)
	Delete(context.Context, string) error
	Search(context.Context, map[string]interface{}) ([]*Property, error)
}

type PropertyAction struct {
	repo PropertyRepo
	log  *log.Helper
}

func NewPropertyAction(repo PropertyRepo, logger log.Logger) *PropertyAction {
	return &PropertyAction{repo: repo, log: log.NewHelper(logger)}
}

func (uc *PropertyAction) GetProperty(ctx context.Context, id string) (*Property, error) {
	uc.log.WithContext(ctx).Infof("GetProperty: %s", id)
	property, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return property, nil
}

func (uc *PropertyAction) ListLodging(ctx context.Context, filters map[string]interface{}) ([]*Property, error) {
	var filtersMask map[string]interface{}
	if filters["userId"] != nil && filters["userId"] != "" {
		filtersMask["userId"] = filters["userId"]
	}

	properties, err := uc.repo.Search(ctx, filtersMask)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (uc *PropertyAction) SearchProperties(ctx context.Context, filters map[string]interface{}) ([]*Property, error) {
	uc.log.WithContext(ctx).Infof("SearchProperties: %s", filters)
	properties, err := uc.repo.Search(ctx, filters)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (uc *PropertyAction) CreateProperty(ctx context.Context, p *Property) (*Property, error) {
	uc.log.WithContext(ctx).Infof("CreateProperty: %s", p.Title)
	res, err := uc.repo.Save(ctx, p)
	if err != nil {
		fmt.Println("error creating property: ", err)
	}
	fmt.Println("create property result: ", res)
	return res, err
}

func (uc *PropertyAction) UpdateProperty(ctx context.Context, p *Property) (*Property, error) {
	uc.log.WithContext(ctx).Infof("UpdateProperty: %s", p.Title)
	res, err := uc.repo.Update(ctx, p)
	if err != nil {
		fmt.Println("error updating property: ", err)
	}
	fmt.Println("update property result: ", res)
	return res, err
}

func (uc *PropertyAction) DeleteProperty(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("Delete Property: %s", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *PropertyAction) GetRealtorStats(ctx context.Context, userId string) (map[string]int64, error) {
	return uc.repo.Stats(ctx, userId)
}
