package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"radioking-app/internal/api/http/authentication"
	"radioking-app/internal/api/http/handlers"
	"radioking-app/internal/config"
	"radioking-app/internal/domain/services"
	"radioking-app/internal/infrastructure/db"
	"radioking-app/internal/infrastructure/messaging"
	"radioking-app/internal/infrastructure/repositories"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func main() {
	cfg := loadConfiguration()

	dbInstance, err := initDb()

	// Initialize messaging
	publisher, consumer, err := initMessaging(cfg)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
	defer consumer.Close()

	// Initialize repositories
	playlistRepo := repositories.NewPlaylistRepository(dbInstance)
	trackPlayRepo := repositories.NewTrackPlayRepository(dbInstance)

	// Initialize services
	playlistService := &services.PlaylistService{Repo: playlistRepo}
	playlistPlayService := services.NewPlaylistPlayService(playlistService, publisher)
	trackPlayService := services.NewTrackPlayService(trackPlayRepo)

	// Initialize application service
	playlistApplicationService := services.NewPlaylistApplicationService(playlistService, playlistPlayService)

	// Initialize consumer service
	consumerService := services.NewTrackPlayConsumerService(consumer, trackPlayService)

	// Start consumer service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumerService.Start(ctx); err != nil {
		log.Printf("Failed to start consumer service: %v", err)
	}
	defer consumerService.Stop()

	// Initialize HTTP router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	if cfg.Auth.Enabled {
		jwtMiddleware := initAuthMiddleware(cfg)
		router.Use(jwtMiddleware.Middleware())
	}

	// Initialize handlers
	handler := handlers.NewPlaylistHandler(playlistService, playlistApplicationService)
	handler.Routes(router)

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	cancel()

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server shutdown complete")
}

func initDb() (*gorm.DB, error) {
	dbInstance, err := db.InitDb()
	if err != nil {
		panic(err)
	}
	return dbInstance, err
}

func initAuthMiddleware(cfg *config.Config) *authentication.JWTMiddleware {
	jwtMiddleware := authentication.NewJWTMiddleware(cfg.Auth.KeycloakURL, cfg.Auth.Realm)
	if err := jwtMiddleware.LoadPublicKeys(); err != nil {
		panic(fmt.Errorf("failed to load public keys: %w", err))
	}
	return jwtMiddleware
}

func loadConfiguration() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}
	return cfg
}

func initMessaging(cfg *config.Config) (messaging.MessagePublisher, messaging.MessageConsumer, error) {
	if err := messaging.InitRabbitMQInfrastructure(cfg.Messaging.RabbitMQ); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize RabbitMQ infrastructure: %w", err)
	}

	publisher, err := messaging.NewRabbitMQPublisher(cfg.Messaging.RabbitMQ)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize publisher: %w", err)
	}

	consumer, err := messaging.NewRabbitMQConsumer(cfg.Messaging.RabbitMQ)
	if err != nil {
		publisher.Close()
		return nil, nil, fmt.Errorf("failed to initialize consumer: %w", err)
	}

	return publisher, consumer, nil
}
