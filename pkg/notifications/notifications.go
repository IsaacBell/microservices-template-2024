package notifications

import (
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	"microservices-template-2024/pkg/cache"
	consultants_biz "microservices-template-2024/pkg/consultants/biz"
	notifications_biz "microservices-template-2024/pkg/notifications/biz"
	notifications_clients "microservices-template-2024/pkg/notifications/clients"
	"microservices-template-2024/pkg/users"
	"os"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	client      *notifications_clients.CourierClient
	ProviderSet = wire.NewSet(Notify, NotifyCommunicationSent, NotifyError, SendSystemAlert)
)

type Notification struct {
	gorm.Model
	UserId   string
	Data     *notifications_biz.NotificationData
	Metadata *notifications_biz.NotificationMetadata
}

func setClient() {
	if client == nil {
		client = notifications_clients.NewCourierClient()
	}
}

func ChangeClient(c *notifications_clients.CourierClient) {
	client = c
}

func setMetadata(notif *Notification) {
	if notif.Metadata.IsError {
		notif.Metadata.Priority = "critical"
	} else if notif.Metadata.IsWarning {
		notif.Metadata.Priority = "medium"
	} else {
		notif.Metadata.Priority = "default"
	}

	if notif.UserId == "" {
		notif.UserId = notif.Data.Recipient.Id
	}
}

func Notify(notif *Notification) error {
	setClient()
	setMetadata(notif)

	err := client.SendNotification(notif.Data, notif.Metadata)
	if err != nil {
		return err
	}

	return nil
}

func NotifyCommunicationSent(from, fromId, toId, msg string) {
	var recip *v1.User
	var err error

	cachedUser := users.UserFromCache(toId)
	if cachedUser == nil {
		recip, err = users.Get(toId, true)
		cache.CacheRecord("user", users.UserCacheKey(toId), toId, recip)
	} else {
		recip = biz.UserToProtoData(cachedUser)
	}
	if err != nil {
		return
	}
	Notify(&Notification{
		UserId: toId,
		Data: &notifications_biz.NotificationData{
			Msg:      msg,
			From:     from,
			CommType: consultants_biz.COMM_TYPE_FromClient.String(),
			Recipient: &notifications_biz.Recipient{
				Id:        toId,
				SenderId:  fromId,
				Email:     recip.Email,
				Phone:     recip.PhoneNumber,
				FirstName: recip.FirstName,
				LastName:  recip.LastName,
			},
		},
		Metadata: notifications_biz.DefaultNotifMetadata(),
	})
}

func NotifyError(uid string, reportedErr error) error {
	u, err := users.Get(uid, true)
	if err != nil {
		return err
	}
	err = Notify(&Notification{
		UserId: uid,
		Data: &notifications_biz.NotificationData{
			From:     "system",
			Msg:      `There was a system error. Please contact support. Error Message: ` + reportedErr.Error(),
			CommType: consultants_biz.COMM_TYPE_FromSystem.String(),
			Recipient: &notifications_biz.Recipient{
				Id:        uid,
				SenderId:  os.Getenv("SYSTEM_SENDER_ID"),
				Email:     u.Email,
				Phone:     u.PhoneNumber,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			},
			Options: map[string]bool{},
		},
		Metadata: notifications_biz.SystemNotifMetadata(),
	})

	if err != nil {
		return err
	}

	return nil
}

func SendSystemAlert(uid, msg string, options map[string]bool) error {
	u, err := users.Get(uid, true)
	if err != nil {
		return err
	}
	err = Notify(&Notification{
		UserId: uid,
		Data: &notifications_biz.NotificationData{
			From:     "system",
			Msg:      msg,
			CommType: consultants_biz.COMM_TYPE_FromSystem.String(),
			Recipient: &notifications_biz.Recipient{
				Id:        uid,
				SenderId:  os.Getenv("SYSTEM_SENDER_ID"),
				Email:     u.Email,
				Phone:     u.PhoneNumber,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			},
			Options: options,
		},
		Metadata: notifications_biz.SystemNotifMetadata(),
	})

	if err != nil {
		return err
	}

	return nil
}
