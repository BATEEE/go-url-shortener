package domain

import "time"

type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Link struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"index;not null" json:"user_id"`
	ShortCode   string    `gorm:"size:32;uniqueIndex;not null" json:"short_code"`
	OriginalUrl string    `gorm:"not null;index" json:"original_url"`
	Clicks      int       `gorm:"default:0" json:"clicks"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
