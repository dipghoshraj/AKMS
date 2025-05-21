package store

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DataBase *gorm.DB

func InitDB() {

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_USER", "dip"),
		GetEnv("DB_PASSWORD", "password"),
		GetEnv("DB_NAME", "servic_db_2"),
		GetEnv("DB_PORT", "5434"),
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get raw DB object: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	DataBase = db
	fmt.Println("Database connection established successfully.")
}
