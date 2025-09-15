package services

import "radioking-app/internal/domain/models"

type IPlaylistService interface {
	CreatePlaylist(playlist *models.Playlist) error
	ListPlaylists() ([]*models.Playlist, error)
	GetPlaylist(id int) (*models.Playlist, error)
}
