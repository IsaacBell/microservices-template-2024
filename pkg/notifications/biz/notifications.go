package notifications_biz

import (
	"time"
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
	From      string
	CommType  string
	Msg       string
	Recipient *Recipient
	Options   map[string]bool
	// Channel        []string
}

type NotificationMetadata struct {
	Priority      string
	IsWarning     bool
	IsError       bool
	IsSystemAlert bool
	WillExpire    bool
	ExpiresAt     time.Time
	SentAt        time.Time
}

func DefaultNotifMetadata() *NotificationMetadata {
	return &NotificationMetadata{
		Priority:      "default",
		WillExpire:    false,
		IsWarning:     false,
		IsError:       false,
		IsSystemAlert: false,
		SentAt:        time.Now(),
	}
}

func WarningNotifMetadata() *NotificationMetadata {
	return &NotificationMetadata{
		Priority:      "medium",
		WillExpire:    false,
		IsWarning:     true,
		IsError:       false,
		IsSystemAlert: false,
		SentAt:        time.Now(),
	}
}

func ErrorNotifMetadata() *NotificationMetadata {
	return &NotificationMetadata{
		Priority:      "medium",
		WillExpire:    false,
		IsWarning:     false,
		IsError:       true,
		IsSystemAlert: false,
		SentAt:        time.Now(),
	}
}

func SystemNotifMetadata() *NotificationMetadata {
	return &NotificationMetadata{
		Priority:      "default",
		WillExpire:    false,
		IsWarning:     false,
		IsError:       false,
		IsSystemAlert: true,
		SentAt:        time.Now(),
	}
}

func NotifMetadataWithExpiry(expiry time.Time) *NotificationMetadata {
	return &NotificationMetadata{
		Priority:      "default",
		WillExpire:    true,
		IsWarning:     false,
		IsError:       false,
		IsSystemAlert: false,
		SentAt:        time.Now(),
		ExpiresAt:     expiry,
	}
}
