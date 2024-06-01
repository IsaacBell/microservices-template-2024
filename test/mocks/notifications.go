package mocks

import (
	"github.com/stretchr/testify/mock"
	courier "github.com/trycourier/courier-go/v2"
)

type MockCourierClient struct {
	courier.Client
	// notifications_clients.CourierClient
	mock.Mock
}

// func (m *MockCourierClient) connect() error {
// 	return nil
// }

// func (m *MockCourierClient) SendNotification(data *notifications_biz.NotificationData, metadata *notifications_biz.NotificationMetadata) error {
// 	m.Called(data, metadata)
// 	return nil
// }
