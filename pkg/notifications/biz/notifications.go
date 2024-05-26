package notifications_biz

import (
	"time"

	"gorm.io/gorm"
)

type Recipient struct {
	Id        string
	SenderId  string
	Email     string
	Phone     string
	FirstName string
	LastName  string
}

type NotificationData struct {
	gorm.Model
	From      string
	CommType  string
	Msg       string
	Recipient *Recipient
	Options   map[string]bool
	// Channel        []string
}

type NotificationMetadata struct {
	gorm.Model
	Priority      string
	IsWarning     bool
	IsError       bool
	IsSystemAlert bool
	WillExpire    bool
	ExpiresAt     time.Time
	SentAt        time.Time
}
