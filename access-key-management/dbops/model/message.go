package model

import "time"

type KafkaMessage struct {
	HashKey         string    `json:"hashkey"`
	RateLimitPerMin int64     `json:"rate_limit_per_min"`
	ExpiresAt       time.Time `json:"expires_at"`
	Disabled        bool      `json:"disabled"`
}
