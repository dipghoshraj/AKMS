package model

import "time"

type Token struct {
	Hashkey            string    `gorm:"primaryKey;size:255;not null" json:"hashkey"`
	ExpiresAt          time.Time `gorm:"not null" json:"expires_at"`
	Disabled           bool      `gorm:"not null" json:"disabled"`
	RateLimitPerMinute int64     `gorm:"not null" json:"rate_limit_per_min"`
}

type RedisMeta struct {
	Disabled  bool      `json:"disabled"`
	ExpiresAt time.Time `json:"expires_at"`
	RateLimit int64     `json:"rate_limit_per_min"`
}
