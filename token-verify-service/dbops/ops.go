package dbops

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type tokenOps struct {
	db    *gorm.DB
	redis *redis.Client
}

type OpsManager struct {
	TokenOps TokenOps
}

func NewOpsManager(db *gorm.DB, redis *redis.Client) *OpsManager {
	return &OpsManager{
		TokenOps: &tokenOps{db: db, redis: redis},
	}
}
