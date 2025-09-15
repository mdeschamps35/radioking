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

func (s *PlaylistPlayService) PlayPlaylist(playlistID int) error {
	// 1. Récupérer la playlist avec ses tracks
	playlist, err := s.playlistService.GetPlaylist(playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist %d: %w", playlistID, err)
	}

	if len(playlist.Tracks) == 0 {
		log.Printf("Playlist %d is empty, nothing to play", playlistID)
		return nil
	}

	log.Printf("Starting to play playlist %d with %d tracks", playlistID, len(playlist.Tracks))

	// 2. Pour chaque track, publier un événement TrackPlayed
	playedAt := time.Now()
	for position, track := range playlist.Tracks {
		event := models.TrackPlayedEvent{
			PlaylistID: playlistID,
			TrackID:    int(track.ID),
			TrackTitle: track.Title,
			Artist:     track.Artist,
			Position:   position, // 0-based position
			PlayedAt:   playedAt,
			EventID:    uuid.New().String(),
		}

		err := s.messagePublisher.PublishTrackPlayedEvent(event)
		if err != nil {
			log.Printf("Failed to publish event for track %d (position %d): %v", track.ID, position, err)
			return fmt.Errorf("failed to publish event for track %d: %w", track.ID, err)
		}

		log.Printf("Published event for track '%s' by '%s' at position %d", track.Title, track.Artist, position)
	}

	log.Printf("Successfully published %d track events for playlist %d", len(playlist.Tracks), playlistID)
	return nil
}
