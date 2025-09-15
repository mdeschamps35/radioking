package repositories

import "radioking-app/internal/domain/models"

type IPlaylistRepository interface {
	Create(playlist *models.Playlist) error
	GetAll() ([]*models.Playlist, error)
	GetByID(id int) (*models.Playlist, error)
}
