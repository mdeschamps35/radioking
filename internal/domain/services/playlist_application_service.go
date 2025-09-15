package services

import (
	"fmt"
)

type PlayPlaylistResult struct {
	PlaylistID  int
	TracksCount int
	Message     string
}

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

	playlist, err := s.playlistService.GetPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist %d: %w", playlistID, err)
	}

	err = s.playlistPlayService.PlayPlaylist(*playlist)
	if err != nil {
		return nil, fmt.Errorf("failed to play playlist %d: %w", playlistID, err)
	}

	return &PlayPlaylistResult{
		PlaylistID:  playlistID,
		TracksCount: len(playlist.Tracks),
		Message:     "Playlist is being played",
	}, nil
}
