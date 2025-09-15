package models

import "time"

// TrackPlayedEvent représente l'événement envoyé dans RabbitMQ quand une track est jouée
// Cette structure est destinée à la sérialisation JSON pour le messaging
type TrackPlayedEvent struct {
	PlaylistID int64     `json:"playlist_id"`
	TrackID    int64     `json:"track_id"`
	TrackTitle string    `json:"track_title"`
	Artist     string    `json:"artist"`
	Position   int       `json:"position"` // Position dans la playlist (0-based)
	PlayedAt   time.Time `json:"played_at"`
	EventID    string    `json:"event_id"` // UUID unique pour l'événement
}
