package notifications_clients

import (
	"context"
	notifications_biz "core/pkg/notifications/biz"
	"fmt"
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

func connect() *courier.Client {
	onceConnClient.Do(func() {
		courierClient = courier.CreateClient(os.Getenv("COURIER_AUTH_TOKEN"), nil)
	})
	return courierClient
}

func NewCourierClient() *CourierClient {
	return &CourierClient{client: connect()}
}

func NewCourierTestClient() *CourierClient {
	return &CourierClient{client: &courier.Client{}}
}

func (c *CourierClient) SendNotification(
	data *notifications_biz.NotificationData,
	metadata *notifications_biz.NotificationMetadata,
) error {
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

				// "metadata": metadata,

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
