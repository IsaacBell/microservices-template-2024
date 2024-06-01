package notifications_clients

import notifications_biz "microservices-template-2024/pkg/notifications/biz"

type Client interface {
	SendNotification(data *notifications_biz.NotificationData, metadata *notifications_biz.NotificationMetadata) error
}
