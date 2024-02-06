package entity

import "time"

type User struct {
	ID        string     `gorm:"id"`
	Username  string     `gorm:"username"`
	Password  string     `gorm:"password"`
	Email     string     `gorm:"email"`
	Profile   string     `gorm:"profile"`
	IsEnabled bool       `gorm:"is_enabled"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at"`
}

type Profile struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}
