package main

import (
	"akm/dbops/model"
	"akm/store"
	"fmt"
	"log"
)

func Migrate() {
	// Perform the migration
	err := store.DataBase.AutoMigrate(&model.Token{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully.")
}
