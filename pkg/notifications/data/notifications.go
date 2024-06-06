package notifications_data

import (
	"core/internal/server"
)

func SaveNotification(notif *interface{}) error {
	if err := server.DB.Create(&notif).Error; err != nil {
		return err
	}

	return nil
}
