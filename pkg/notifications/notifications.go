package notifications

import (
	"core/internal/biz"
	"core/pkg/cache"
	consultants_biz "core/pkg/consultants/biz"
	notifications_biz "core/pkg/notifications/biz"
	notifications_clients "core/pkg/notifications/clients"
	"core/pkg/stream"
	"core/pkg/users"
	"os"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	testMode    bool
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
	if client == nil && testMode {
		client = notifications_clients.NewCourierTestClient()
	}
	if client == nil && !testMode {
		client = notifications_clients.NewCourierClient()
	}
}

func ChangeClient(c *notifications_clients.CourierClient) {
	client = c
}

// for failover simulation
func DropClient() {
	client = nil
}

func EnableTestMode() {
	testMode = true
}

func DisableTestMode() {
	testMode = false
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

	if !testMode {
		stream.ProduceKafkaMessage("notifications", notif.Data.Msg)
	}

	err := client.SendNotification(notif.Data, notif.Metadata)
	if err != nil {
		return err
	}

	return nil
}

func getUserWithCache(
	id string,
	failIfNotFound bool,
) (*biz.User, error) {
	cachedUser := users.UserFromCache(id)
	if cachedUser == nil {
		u, err := users.Get(id, true)
		if err != nil {
			return nil, err
		}
		cache.CacheRecord("user", users.UserCacheKey(id), id, u)
		return biz.ProtoToUserData(u), nil
	} else {
		return cachedUser, nil
	}

	return nil, nil
}

func NotifyCommunicationSent(from, fromId, toId, msg string) {
	recip, err := getUserWithCache(toId, false)
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
	u, err := getUserWithCache(uid, false)
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
	u, err := getUserWithCache(uid, false)
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
