package messaging

import (
	"context"
	"radioking-app/internal/domain/models"
)

// MessageConsumer interface for consuming messages
type MessageConsumer interface {
	ConsumeTrackPlayedEvents(ctx context.Context, handler func(models.TrackPlayedEvent) error) error
	Close() error
}
