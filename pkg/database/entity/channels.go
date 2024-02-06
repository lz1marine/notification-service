package entity

type Channel struct {
	ID        string `gorm:"id"`
	Name      string `gorm:"name"`
	IsEnabled bool   `gorm:"is_enabled"`
}

type UserChannel struct {
	ID        string   `gorm:"id"`
	UserID    string   `gorm:"user_id"`
	ChannelID string   `gorm:"channel_id"`
	IsEnabled bool     `gorm:"is_enabled"`
	User      *User    `gorm:"foreignKey:UserID"`
	Channel   *Channel `gorm:"foreignKey:ChannelID"`
}
