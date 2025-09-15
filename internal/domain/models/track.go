package models

import (
	"time"
)

type Track struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	PlaylistID int64  `gorm:"index;not null"`
	Title      string `gorm:"size:255;not null"`
	Artist     string `gorm:"size:255;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
