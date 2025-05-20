package model

import "time"

type Token struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Hashkey            string    `gorm:"size:255;not null" json:"hashkey"`
	ExpiresAt          time.Time `gorm:"not null" json:"expires_at"`
	Disabled           bool      `gorm:"not null" json:"disabled"`
	RateLimitPerMinute int64     `gorm:"not null" json:"rate_limit_per_min"`
}
