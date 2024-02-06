package entity

import "time"

// Message state
const (
	MessagePending = iota
	MessageActive
	MessageSent
	MessageFailed
	MessageStale
	MessageDeleted
)

type Message struct {
	ID         string     `gorm:"id"`
	EventID    string     `gorm:"event_id"`
	Subject    *string    `gorm:"subject"`
	Message    string     `gorm:"message"`
	TemplateID *string    `gorm:"template_id"`
	Status     int        `gorm:"status"`
	ChannelID  string     `gorm:"channel_id"`
	CreatedAt  time.Time  `gorm:"created_at"`
	UpdatedAt  time.Time  `gorm:"updated_at"`
	DeletedAt  *time.Time `gorm:"deleted_at"`
	Version    int        `gorm:"version"`
}
