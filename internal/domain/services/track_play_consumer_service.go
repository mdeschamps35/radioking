package services

import (
	"context"
	"fmt"
	"log"
	"radioking-app/internal/domain/models"
	"radioking-app/internal/infrastructure/messaging"
	"sync"
)

type TrackPlayConsumerService struct {
	consumer     messaging.MessageConsumer
	trackPlaySvc ITrackPlayService
	stopChan     chan struct{}
	wg           sync.WaitGroup
	isRunning    bool
	mu           sync.Mutex
}

func NewTrackPlayConsumerService(consumer messaging.MessageConsumer, trackPlaySvc ITrackPlayService) *TrackPlayConsumerService {
	return &TrackPlayConsumerService{
		consumer:     consumer,
		trackPlaySvc: trackPlaySvc,
		stopChan:     make(chan struct{}),
	}
}

func (s *TrackPlayConsumerService) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("consumer service is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	// Create handler function for processing events
	handler := func(event models.TrackPlayedEvent) error {
		return s.trackPlaySvc.RecordTrackPlay(event)
	}

	err := s.consumer.ConsumeTrackPlayedEvents(ctx, handler)
	if err != nil {
		s.mu.Lock()
		s.isRunning = false
		s.mu.Unlock()
		return fmt.Errorf("failed to start consuming events: %w", err)
	}

	log.Println("Track play consumer service started")

	// Wait for context cancellation or stop signal
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
		}()

		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping track play consumer service")
		case <-s.stopChan:
			log.Println("Stop signal received, stopping track play consumer service")
		}
	}()

	return nil
}

func (s *TrackPlayConsumerService) Stop() {
	s.mu.Lock()
	if !s.isRunning {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	close(s.stopChan)
	s.wg.Wait()

	if s.consumer != nil {
		s.consumer.Close()
	}

	log.Println("Track play consumer service stopped")
}
