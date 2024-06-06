package notifications_clients

import notifications_biz "core/pkg/notifications/biz"

type Client interface {
	SendNotification(data *notifications_biz.NotificationData, metadata *notifications_biz.NotificationMetadata) error
}
