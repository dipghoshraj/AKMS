package model

import "time"

type Token struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Hashkey   string    `gorm:"size:255;not null" json:"hashkey"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Disabled  bool      `gorm:"not null" json:"disabled"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
