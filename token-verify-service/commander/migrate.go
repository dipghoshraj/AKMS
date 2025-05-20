package main

import (
	"fmt"
	"log"
	"tvs/dbops/model"
	"tvs/store"
)

func Migrate() {
	// Perform the migration
	err := store.DataBase.AutoMigrate(&model.Token{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully.")
}
