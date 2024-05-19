package lodging_data

import (
	"context"
	"microservices-template-2024/internal/server"
	lodging_biz "microservices-template-2024/pkg/lodging/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type propertyRepo struct {
	data *Data
	log  *log.Helper
}

func NewPropertyRepo(data *Data, logger log.Logger) lodging_biz.PropertyRepo {
	return &propertyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *propertyRepo) Get(ctx context.Context, id string) (*lodging_biz.Property, error) {
	var property *lodging_biz.Property
	err := server.DB.Scopes(server.Active).First(&property, id).Error
	if err != nil {
		return nil, err
	}

	return property, nil
}

func (r *propertyRepo) Save(ctx context.Context, p *lodging_biz.Property) (*lodging_biz.Property, error) {
	if p.ID != "" {
		if err := server.DB.Save(&p).Error; err != nil {
			return nil, err
		} else {
			return p, nil
		}
	}

	if err := server.DB.Omit("ID").FirstOrCreate(&p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (r *propertyRepo) Update(ctx context.Context, property *lodging_biz.Property) (*lodging_biz.Property, error) {
	if err := server.DB.Save(&property).Error; err != nil {
		return nil, err
	}
	return property, nil
}

func (r *propertyRepo) Stats(ctx context.Context, userId string) (map[string]int64, error) {
	var result struct {
		TotalProperties int64 `gorm:"column:total_properties"`
		AveragePrice    int64 `gorm:"column:average_price"`
		TotalPrice      int64 `gorm:"column:total_price"`
	}

	if err := server.DB.
		Model(&lodging_biz.Property{}).
		Where("userId = ?", userId).
		Group("address").
		Select("COUNT(*) AS total_properties, AVG(price) AS average_price, SUM(price) AS total_price").
		Scan(&result).Error; err != nil {
		return nil, err
	}

	return map[string]int64{
		"total_properties": result.TotalProperties,
		"average_price":    result.AveragePrice,
		"total_price":      result.TotalPrice,
	}, nil
}

func (r *propertyRepo) Delete(ctx context.Context, id string) error {
	var property *lodging_biz.Property
	if err := server.DB.Scopes(server.Active).First(&property, id).Error; err != nil {
		return err
	}
	property.Deleted = true

	if err := server.DB.Save(&property).Error; err != nil {
		return err
	}
	return nil
}

func (r *propertyRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*lodging_biz.Property, error) {
	var properties []*lodging_biz.Property
	query := server.DB.Scopes(server.Active)

	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if sold, ok := filters["sold"]; ok {
		query = query.Where("sold = ?", sold)
	}
	if rooms, ok := filters["rooms"]; ok {
		query = query.Where("rooms = ?", rooms)
	}
	if priceGte, ok := filters["price_gte"]; ok {
		query = query.Where("price >= ?", priceGte)
	}
	if priceLte, ok := filters["price_lte"]; ok {
		query = query.Where("price <= ?", priceLte)
	}

	err := query.Find(&properties).Error
	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *propertyRepo) Within(ctx context.Context, latitude, longitude float64, distanceInMile int) ([]*lodging_biz.Property, error) {
	var properties []*lodging_biz.Property
	distanceInMeter := float64(distanceInMile) * 1609.34 // approx

	err := server.DB.Scopes(server.Active).Where("ST_Distance(location, ST_MakePoint(?, ?)) < ?", longitude, latitude, distanceInMeter).Find(&properties).Error
	if err != nil {
		return nil, err
	}

	return properties, nil
}
