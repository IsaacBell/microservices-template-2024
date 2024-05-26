package consultants_data

import (
	"context"
	"encoding/json"
	"fmt"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/internal/server"
	"microservices-template-2024/pkg/cache"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	"microservices-template-2024/pkg/notifications"
	notifications_biz "microservices-template-2024/pkg/notifications/biz"

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
	err := server.DB.First(&consultant, id).Error
	if err != nil {
		return nil, err
	}

	return consultant, nil
}

func notifyCommunicationSent(ctx context.Context, comm *consultants_biz.Communication) error {
	var recip *biz.User
	cacheKey := "user:" + comm.RecipientID
	if cached, err := cache.Cache(ctx).Get(cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &recip); err != nil {
			fmt.Println("cache miss: ", err)
			
			err = server.DB.First(&recip, comm.RecipientID).Error
			if err != nil {
				fmt.Println("error retrieving user from DB: ", err)
				
				return err
			}

			cache.CacheRecord("user", cacheKey, comm.RecipientID, recip)
		}
	}

	notif := notifications.Notification{
		UserId: comm.RecipientID,
		Data: &notifications_biz.NotificationData{
			CommType: comm.CommType.String(),
			Msg:      comm.Msg,
			Options:  comm.Options,
			From:     comm.From,
			Recipient: &notifications_biz.Recipient{
				Id:        comm.RecipientID,
				SenderId:  comm.UserID,
				Email:     recip.Email,
				Phone:     recip.PhoneNumber,
				FirstName: recip.FirstName,
				LastName:  recip.LastName,
			},
		},
		Metadata: &notifications_biz.NotificationMetadata{
			Priority:   "default",
			WillExpire: false,
		},
	}

	err := notifications.Notify(&notif)
	if err != nil {
		fmt.Println("error sending notification: ", err)
		return err
	}

	return nil
}

func (r *consultantRepo) SaveCommunication(ctx context.Context, comm *consultants_biz.Communication) (*consultants_biz.Communication, error) {
	go notifyCommunicationSent(ctx, comm)
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

	if err := server.DB.Omit("ID").FirstOrCreate(&c).Error; err != nil {
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
