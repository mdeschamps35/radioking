package messaging

import "radioking-app/internal/domain/models"

type MessagePublisher interface {
	PublishTrackPlayedEvent(event models.TrackPlayedEvent) error
	Close() error
}
