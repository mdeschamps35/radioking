package repositories

import (
	"fmt"
	"radioking-app/internal/domain/models"

	"gorm.io/gorm"
)

type PlaylistRepository struct {
	DB *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) *PlaylistRepository {
	return &PlaylistRepository{DB: db}
}

func (r *PlaylistRepository) Create(playlist *models.Playlist) error {
	if err := r.DB.Create(playlist).Error; err != nil {
		return fmt.Errorf("failed to create playlist in database: %w", err)
	}
	return nil
}

func (r *PlaylistRepository) GetAll() ([]*models.Playlist, error) {
	var playlists []*models.Playlist
	if err := r.DB.Preload("Tracks").Find(&playlists).Error; err != nil {
		return nil, fmt.Errorf("failed to get playlists from database: %w", err)
	}
	return playlists, nil
}

func (r *PlaylistRepository) GetByID(id int) (*models.Playlist, error) {
	var playlist models.Playlist
	err := r.DB.Preload("Tracks").First(&playlist, id).Error

	if err != nil {
		return nil, err
	}

	return &playlist, nil
}
