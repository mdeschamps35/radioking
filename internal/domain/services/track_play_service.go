package services

import (
	"fmt"
	"log"
	"radioking-app/internal/domain/models"
)

type TrackPlayRepository interface {
	Create(trackPlay *models.TrackPlay) error
	GetByPlaylistID(playlistID int) ([]*models.TrackPlay, error)
	GetByTrackID(trackID int) ([]*models.TrackPlay, error)
}

type TrackPlayService struct {
	repository TrackPlayRepository
}

func NewTrackPlayService(repository TrackPlayRepository) *TrackPlayService {
	return &TrackPlayService{
		repository: repository,
	}
}

func (s *TrackPlayService) RecordTrackPlay(event models.TrackPlayedEvent) error {
	trackPlay := &models.TrackPlay{
		PlaylistID: event.PlaylistID,
		TrackID:    event.TrackID,
		Position:   event.Position,
		PlayedAt:   event.PlayedAt,
	}

	if err := s.repository.Create(trackPlay); err != nil {
		return fmt.Errorf("failed to record track play: %w", err)
	}

	log.Printf("Recorded track play: PlaylistID=%d, TrackID=%d, Position=%d, PlayedAt=%v",
		trackPlay.PlaylistID, trackPlay.TrackID, trackPlay.Position, trackPlay.PlayedAt)

	return nil
}

func (s *TrackPlayService) GetPlaylistPlays(playlistID int) ([]*models.TrackPlay, error) {
	plays, err := s.repository.GetByPlaylistID(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist plays: %w", err)
	}
	return plays, nil
}

func (s *TrackPlayService) GetTrackPlays(trackID int) ([]*models.TrackPlay, error) {
	plays, err := s.repository.GetByTrackID(trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track plays: %w", err)
	}
	return plays, nil
}
