package repositories

import (
	"radioking-app/internal/domain/models"

	"gorm.io/gorm"
)

type TrackPlayRepository struct {
	DB *gorm.DB
}

func NewTrackPlayRepository(db *gorm.DB) *TrackPlayRepository {
	return &TrackPlayRepository{DB: db}
}

func (r *TrackPlayRepository) Create(trackPlay *models.TrackPlay) error {
	return r.DB.Create(trackPlay).Error
}

func (r *TrackPlayRepository) GetByPlaylistID(playlistID int) ([]*models.TrackPlay, error) {
	var trackPlays []*models.TrackPlay
	err := r.DB.Where("playlist_id = ?", playlistID).
		Order("played_at DESC").
		Find(&trackPlays).Error

	if err != nil {
		return nil, err
	}

	return trackPlays, nil
}

func (r *TrackPlayRepository) GetByTrackID(trackID int) ([]*models.TrackPlay, error) {
	var trackPlays []*models.TrackPlay
	err := r.DB.Where("track_id = ?", trackID).
		Order("played_at DESC").
		Find(&trackPlays).Error

	if err != nil {
		return nil, err
	}

	return trackPlays, nil
}
