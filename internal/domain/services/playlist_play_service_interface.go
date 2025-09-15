package services

import "radioking-app/internal/domain/models"

type IPlaylistPlayService interface {
	PlayPlaylist(playlist models.Playlist) error
}
