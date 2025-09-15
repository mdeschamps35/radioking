package db

import (
	"fmt"
	"radioking-app/internal/domain/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.Playlist{}, &models.Track{}, &models.TrackPlay{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}
