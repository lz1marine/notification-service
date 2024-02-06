package entity

import "time"

type Topic struct {
	ID        string     `gorm:"id"`
	Name      string     `gorm:"name"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at"`
}

type UserTopic struct {
	ID      string `gorm:"id"`
	UserID  string `gorm:"user_id"`
	TopicID string `gorm:"topic_id"`
	User    *User  `gorm:"foreignKey:UserID" json:"user"`
}
