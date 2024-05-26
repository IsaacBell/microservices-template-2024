package notifications_data

import (
	"microservices-template-2024/internal/server"
)

func SaveNotification(notif *interface{}) error {
	if err := server.DB.Create(&notif).Error; err != nil {
		return err
	}

	return nil
}
