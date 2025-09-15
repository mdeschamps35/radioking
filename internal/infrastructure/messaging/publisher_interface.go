package messaging

import "radioking-app/internal/domain/models"

// MessagePublisher interface pour publier des messages
type MessagePublisher interface {
	PublishTrackPlayedEvent(event models.TrackPlayedEvent) error
	Close() error
}
