package models

import (
	"time"
)

type Playlist struct {
	ID        int64   `gorm:"primaryKey;autoIncrement"`
	Name      string  `gorm:"size:255;not null"`
	Tracks    []Track `gorm:"foreignKey:PlaylistID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
