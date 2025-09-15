package services

import "radioking-app/internal/domain/models"

// ITrackPlayService interface for managing track plays
type ITrackPlayService interface {
	RecordTrackPlay(event models.TrackPlayedEvent) error
	GetPlaylistPlays(playlistID int) ([]*models.TrackPlay, error)
	GetTrackPlays(trackID int) ([]*models.TrackPlay, error)
}
