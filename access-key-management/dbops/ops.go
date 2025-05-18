package dbops

import (
	"log"
	"sync"
)

var instance *OpsManager
var once sync.Once

type tokenOps struct{}

type OpsManager struct {
	TokenOps TokenOps
}

func NewAppRepository() TokenOps {
	return &tokenOps{}
}

func InitOpsManager() {
	once.Do(func() {
		instance = &OpsManager{
			TokenOps: &tokenOps{},
		}
		log.Println("RepositoryManager initialized")
	})
}

func GetOpsManager() *OpsManager {
	if instance == nil {
		InitOpsManager()
	}
	return instance
}
