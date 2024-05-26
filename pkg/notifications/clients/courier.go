package notifications_clients

import (
	"context"
	"fmt"
	notifications_biz "microservices-template-2024/pkg/notifications/biz"
	"os"
	"sync"

	courier "github.com/trycourier/courier-go/v2"
)

var (
	onceConnClient sync.Once
	courierClient  *courier.Client
)

type CourierClient struct {
	Client
	client *courier.Client
}

func NewNotificationsClient() *CourierClient {
	return &CourierClient{client: connect()}
}

func connect() *courier.Client {
	onceConnClient.Do(func() {
		courierClient = courier.CreateClient(os.Getenv("COURIER_AUTH_TOKEN"), nil)
	})
	return courierClient
}

func (c *CourierClient) SendNotification(data *notifications_biz.NotificationData) error {
	requestID, err := courierClient.SendMessage(
		context.Background(),
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email": "isaacbell388@gmail.com",
				},
				"template": os.Getenv("COURIER_DEFAULT_MSG_TEMPLATE"),
				"data": map[string]string{
					"from": "isaac.bell@toptal.com",
					"msg":  "msg",
				},
			},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Sent message %s\n", requestID)

	requestID, err = courierClient.SendMessage(
		context.Background(),
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email":        data.Recipient.Id,
					"phone_number": data.Recipient.Phone,
				},
				"template": os.Getenv("COURIER_DEFAULT_MSG_TEMPLATE"),
				"data": map[string]string{
					"from": data.From,
					"msg":  data.Msg,
				},
				// "routing": map[string]interface{}{
				// 	"method": "single",
				// 	"channels": []string{"sms", "email", "inbox"},
				// },
			},
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Sent message %s\n", requestID)
	return nil
}
