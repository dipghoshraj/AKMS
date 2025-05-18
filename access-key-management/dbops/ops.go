package dbops

import (
	"gorm.io/gorm"
)

type tokenOps struct {
	db *gorm.DB
}

type OpsManager struct {
	TokenOps TokenOps
}

func NewOpsManager(db *gorm.DB) *OpsManager {
	return &OpsManager{
		TokenOps: &tokenOps{db: db},
	}
}
