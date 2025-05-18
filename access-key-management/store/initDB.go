package store

import (
	"akm/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DataBase *gorm.DB

func InitDB() {

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "password"),
		config.GetEnv("DB_NAME", "myapp"),
		config.GetEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	DataBase = db
	fmt.Println("Database connection established successfully.")
}
