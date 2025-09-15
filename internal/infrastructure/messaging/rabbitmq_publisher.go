package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"radioking-app/internal/config"
	"radioking-app/internal/domain/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

func NewRabbitMQPublisher(cfg config.RabbitMQConfig) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	publisher := &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}

	log.Printf("RabbitMQ Publisher connected to %s", cfg.URL)
	return publisher, nil
}

func (p *RabbitMQPublisher) PublishTrackPlayedEvent(event models.TrackPlayedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.Publish(
		p.config.Exchange,
		p.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published track played event: PlaylistID=%d, TrackID=%d, Position=%d",
		event.PlaylistID, event.TrackID, event.Position)
	return nil
}

func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}
