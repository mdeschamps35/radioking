package models

import "time"

type TrackPlay struct {
	ID         int64     `gorm:"primaryKey"`
	PlaylistID int64     `gorm:"not null;index"`
	TrackID    int64     `gorm:"not null;index"`
	Position   int       `gorm:"not null"`
	PlayedAt   time.Time `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	// Relations
	Playlist Playlist `gorm:"foreignKey:PlaylistID"`
	Track    Track    `gorm:"foreignKey:TrackID"`
}
