package notifications

import (
	notifications_biz "microservices-template-2024/pkg/notifications/biz"
	notifications_clients "microservices-template-2024/pkg/notifications/clients"
	"time"

	// notifications_data "microservices-template-2024/pkg/notifications/data"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	client      notifications_clients.Client
	ProviderSet = wire.NewSet(Notify)
)

type Notification struct {
	gorm.Model
	UserId   string
	Data     *notifications_biz.NotificationData
	Metadata *notifications_biz.NotificationMetadata
}

func setMetadata(notif *Notification) {
	if notif.Metadata.IsError {
		notif.Metadata.Priority = "critical"
	} else if notif.Metadata.IsWarning {
		notif.Metadata.Priority = "medium"
	} else {
		notif.Metadata.Priority = "default"
	}

	notif.Metadata.SentAt = time.Now()
	notif.UserId = notif.Data.Recipient.Id
}

func Notify(notif *Notification) error {
	if client == nil {
		client = notifications_clients.NewNotificationsClient()
	}

	setMetadata(notif)

	err := client.SendNotification(notif.Data)
	if err != nil {
		return err
	}

	return nil
}

// func (notif *Notification) Save() error {
// 	return notifications_data.SaveNotification(notif.Data)
// }
