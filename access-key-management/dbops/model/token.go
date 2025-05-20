package model

import "time"

type Token struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Hashkey            string    `gorm:"size:255;not null" json:"hashkey"`
	RateLimitPerMinute int64     `gorm:"not null" json:"rate_limit_per_minute"`
	ExpiresAt          time.Time `gorm:"not null" json:"expires_at"`
	Disabled           bool      `gorm:"not null" json:"disabled"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TokenCreateInput struct {
	RateLimitPerMinute int64 `json:"rate_limit_per_minute" binding:"required"`
	ExpiresAt          int   `json:"expires_at" binding:"required"`
}

type TokenUpdateInput struct {
	RateLimitPerMinute int64     `json:"rate_limit_per_minute" binding:"required"`
	ExpiresAt          time.Time `json:"expires_at" binding:"required"`
}
