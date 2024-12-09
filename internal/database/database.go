package database

import (
	"fmt"
	"spageti-x-drone-management/internal/config"
	"spageti-x-drone-management/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.Drone{}, &models.Order{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}
