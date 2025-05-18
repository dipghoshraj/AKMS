package dbops

import (
	"sync"
)

var instance *OpsManager
var once sync.Once

type tokenOps struct{}

type OpsManager struct {
	TokenOps TokenOps
}

func NewOpsManager() *OpsManager {
	return &OpsManager{
		TokenOps: &tokenOps{},
	}
}
