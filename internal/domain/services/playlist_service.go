package services

import (
	"errors"
	"fmt"
	"radioking-app/internal/domain/constants"
	domainErrors "radioking-app/internal/domain/errors"
	"radioking-app/internal/domain/models"
	"radioking-app/internal/infrastructure/repositories"
	"strings"

	"gorm.io/gorm"
)

type PlaylistService struct {
	Repo repositories.IPlaylistRepository
}

func (service *PlaylistService) CreatePlaylist(playlist *models.Playlist) error {
	if err := service.validatePlaylist(playlist); err != nil {
		return err
	}

	if err := service.Repo.Create(playlist); err != nil {
		return domainErrors.NewInternalError("failed to create playlist", err)
	}

	return nil
}

func (service *PlaylistService) validatePlaylist(playlist *models.Playlist) error {
	if strings.TrimSpace(playlist.Name) == "" {
		return domainErrors.ErrEmptyPlaylistName
	}

	if len(playlist.Name) > constants.MaxPlaylistNameLength {
		return domainErrors.NewValidationError(fmt.Sprintf("playlist name too long (max %d characters)", constants.MaxPlaylistNameLength))
	}

	if len(playlist.Tracks) > constants.MaxTracksPerPlaylist {
		return domainErrors.ErrTooManyTracks
	}

	for i, track := range playlist.Tracks {
		if err := service.validateTrack(&track); err != nil {
			return domainErrors.NewValidationError(fmt.Sprintf("track %d invalid: %s", i+1, err.Error()))
		}
	}

	return nil
}

func (service *PlaylistService) validateTrack(track *models.Track) error {
	if strings.TrimSpace(track.Title) == "" {
		return domainErrors.ErrEmptyTrackTitle
	}
	if strings.TrimSpace(track.Artist) == "" {
		return domainErrors.ErrEmptyTrackArtist
	}
	if len(track.Title) > constants.MaxTrackNameLength {
		return domainErrors.NewValidationError(fmt.Sprintf("track title too long (max %d characters)", constants.MaxTrackNameLength))
	}
	if len(track.Artist) > constants.MaxArtistNameLength {
		return domainErrors.NewValidationError(fmt.Sprintf("artist name too long (max %d characters)", constants.MaxArtistNameLength))
	}
	return nil
}

func (service *PlaylistService) ListPlaylists() ([]*models.Playlist, error) {
	playlists, err := service.Repo.GetAll()
	if err != nil {
		return nil, domainErrors.NewInternalError("failed to list playlists", err)
	}
	return playlists, nil
}

func (service *PlaylistService) GetPlaylist(id int) (*models.Playlist, error) {
	if id <= 0 {
		return nil, domainErrors.ErrInvalidPlaylistID
	}

	playlist, err := service.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrPlaylistNotFound
		}
		return nil, domainErrors.NewInternalError("failed to get playlist", err)
	}
	return playlist, nil
}
