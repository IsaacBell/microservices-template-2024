package consultants_data

import (
	"context"
	"fmt"
	"microservices-template-2024/internal/server"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	"microservices-template-2024/pkg/notifications"

	"github.com/go-kratos/kratos/v2/log"
)

type consultantRepo struct {
	data *Data
	log  *log.Helper
}

func NewConsultantRepo(data *Data, logger log.Logger) consultants_biz.ConsultantRepo {
	return &consultantRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *consultantRepo) Get(ctx context.Context, id string) (*consultants_biz.Consultant, error) {
	var consultant *consultants_biz.Consultant
	fmt.Println("here")
	err := server.DB.First(&consultant, id).Error
	if err != nil {
		return nil, err
	}

	return consultant, nil
}

func (r *consultantRepo) SaveCommunication(ctx context.Context, comm *consultants_biz.Communication) (*consultants_biz.Communication, error) {
	go notifications.NotifyCommunicationSent(comm.From, comm.UserID, comm.RecipientID, comm.Msg)
	return comm, nil
}

func (r *consultantRepo) Save(ctx context.Context, c *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	if c.ID != "" {
		if err := server.DB.Save(&c).Error; err != nil {
			return nil, err
		} else {
			return c, nil
		}
	}

	if err := server.DB.Create(&c).Error; err != nil {
		return nil, err
	}

	return c, nil
}

func (r *consultantRepo) Update(ctx context.Context, consultant *consultants_biz.Consultant) (*consultants_biz.Consultant, error) {
	if err := server.DB.Save(&consultant).Error; err != nil {
		return nil, err
	}
	return consultant, nil
}

func (r *consultantRepo) Delete(ctx context.Context, id string) error {
	var consultant *consultants_biz.Consultant
	if err := server.DB.First(&consultant, id).Error; err != nil {
		return err
	}

	if err := server.DB.Delete(&consultant).Error; err != nil {
		return err
	}
	return nil
}

func (r *consultantRepo) Search(ctx context.Context, filters map[string]interface{}) ([]*consultants_biz.Consultant, error) {
	var consultants []*consultants_biz.Consultant
	query := server.DB

	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if specializations, ok := filters["specializations"]; ok {
		query = query.Where("specializations @> ?", specializations)
	}
	if languages, ok := filters["languages"]; ok {
		query = query.Where("languages @> ?", languages)
	}
	if yearsOfExperienceGte, ok := filters["years_of_experience_gte"]; ok {
		query = query.Where("years_of_experience >= ?", yearsOfExperienceGte)
	}
	if yearsOfExperienceLte, ok := filters["years_of_experience_lte"]; ok {
		query = query.Where("years_of_experience <= ?", yearsOfExperienceLte)
	}

	err := query.Find(&consultants).Error
	if err != nil {
		return nil, err
	}

	return consultants, nil
}
