package services

import (
	"fmt"
)

// PlayPlaylistResult contient le résultat de l'opération de lecture
type PlayPlaylistResult struct {
	PlaylistID  int
	TracksCount int
	Message     string
}

// IPlaylistApplicationService interface pour les opérations applicatives sur les playlists
type IPlaylistApplicationService interface {
	PlayPlaylist(playlistID int) (*PlayPlaylistResult, error)
}

type PlaylistApplicationService struct {
	playlistService     IPlaylistService
	playlistPlayService IPlaylistPlayService
}

func NewPlaylistApplicationService(
	playlistService IPlaylistService,
	playlistPlayService IPlaylistPlayService,
) *PlaylistApplicationService {
	return &PlaylistApplicationService{
		playlistService:     playlistService,
		playlistPlayService: playlistPlayService,
	}
}

func (s *PlaylistApplicationService) PlayPlaylist(playlistID int) (*PlayPlaylistResult, error) {
	// 1. Vérifier que la playlist existe et récupérer les infos
	playlist, err := s.playlistService.GetPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist %d: %w", playlistID, err)
	}

	// 2. Jouer la playlist (envoyer les événements)
	err = s.playlistPlayService.PlayPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to play playlist %d: %w", playlistID, err)
	}

	// 3. Retourner le résultat avec les informations nécessaires
	return &PlayPlaylistResult{
		PlaylistID:  playlistID,
		TracksCount: len(playlist.Tracks),
		Message:     "Playlist is being played",
	}, nil
}
