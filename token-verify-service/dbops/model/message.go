package model

import "time"

type KafkaMessage struct {
	HashKey         string    `json:"hashkey"`
	RateLimitPerMin int64     `json:"rate_limit_per_min"`
	ExpiresAt       time.Time `json:"expires_at"`
	Disabled        bool      `json:"disabled"`
	EventType       string    `json:"event_type"`
	ReqID           string    `json:"request_id"`
}

type DisableMessage struct {
	HashKey   string `json:"hashkey"`
	Disabled  bool   `json:"disabled"`
	EventType string `json:"event_type"`
	ReqID     string `json:"request_id"`
}
