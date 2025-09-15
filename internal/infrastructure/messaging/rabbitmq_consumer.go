package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"radioking-app/internal/config"
	"radioking-app/internal/domain/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

func NewRabbitMQConsumer(cfg config.RabbitMQConfig) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	consumer := &RabbitMQConsumer{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}

	log.Printf("RabbitMQ Consumer connected to %s", cfg.URL)
	return consumer, nil
}

func (c *RabbitMQConsumer) ConsumeTrackPlayedEvents(ctx context.Context, handler func(models.TrackPlayedEvent) error) error {
	msgs, err := c.channel.Consume(
		c.config.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping RabbitMQ consumer")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ messages channel closed")
					return
				}

				var event models.TrackPlayedEvent
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Printf("Failed to unmarshal track played event: %v", err)
					msg.Nack(false, false)
					continue
				}

				log.Printf("Consumed track played event: PlaylistID=%d, TrackID=%d, Position=%d",
					event.PlaylistID, event.TrackID, event.Position)

				if err := handler(event); err != nil {
					log.Printf("Failed to process track played event: %v", err)
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
