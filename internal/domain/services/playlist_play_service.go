package services

import (
	"fmt"
	"log"
	"radioking-app/internal/domain/models"
	"radioking-app/internal/infrastructure/messaging"
	"time"

	"github.com/google/uuid"
)

type PlaylistPlayService struct {
	playlistService  IPlaylistService
	messagePublisher messaging.MessagePublisher
}

func NewPlaylistPlayService(playlistService IPlaylistService, publisher messaging.MessagePublisher) *PlaylistPlayService {
	return &PlaylistPlayService{
		playlistService:  playlistService,
		messagePublisher: publisher,
	}
}

func (s *PlaylistPlayService) PlayPlaylist(playlist models.Playlist) error {

	if len(playlist.Tracks) == 0 {
		log.Printf("Playlist %d is empty, nothing to play", playlist.ID)
		return nil
	}

	log.Printf("Starting to play playlist %d with %d tracks", playlist.ID, len(playlist.Tracks))

	err := sendTracksEvents(playlist, s)
	if err != nil {
		return err
	}

	log.Printf("Successfully published %d track events for playlist %d", len(playlist.Tracks), playlist.ID)
	return nil
}

func sendTracksEvents(playlist models.Playlist, s *PlaylistPlayService) error {
	for position, track := range playlist.Tracks {
		event := models.TrackPlayedEvent{
			PlaylistID: playlist.ID,
			TrackID:    track.ID,
			TrackTitle: track.Title,
			Artist:     track.Artist,
			Position:   position, // 0-based position
			PlayedAt:   time.Now(),
			EventID:    uuid.New().String(),
		}

		err := s.messagePublisher.PublishTrackPlayedEvent(event)
		if err != nil {
			log.Printf("Failed to publish event for track %d (position %d): %v", track.ID, position, err)
			return fmt.Errorf("failed to publish event for track %d: %w", track.ID, err)
		}

		log.Printf("Published event for track '%s' by '%s' at position %d", track.Title, track.Artist, position)
	}
	return nil
}
