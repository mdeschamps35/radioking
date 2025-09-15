package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Messaging MessagingConfig `mapstructure:"messaging"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	KeycloakURL string `mapstructure:"keycloak_url"`
	Realm       string `mapstructure:"realm"`
}

type MessagingConfig struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type RabbitMQConfig struct {
	URL        string `mapstructure:"url"`
	Exchange   string `mapstructure:"exchange"`
	Queue      string `mapstructure:"queue"`
	RoutingKey string `mapstructure:"routing_key"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Variables d'environnement avec priorité
	viper.AutomaticEnv()
	viper.SetEnvPrefix("RADIOKING")

	// Mapping des variables d'environnement
	viper.BindEnv("server.port", "RADIOKING_SERVER_PORT")
	viper.BindEnv("auth.enabled", "RADIOKING_AUTH_ENABLED")
	viper.BindEnv("auth.keycloak_url", "RADIOKING_AUTH_KEYCLOAK_URL")
	viper.BindEnv("auth.realm", "RADIOKING_AUTH_REALM")
	viper.BindEnv("messaging.rabbitmq.url", "RADIOKING_RABBITMQ_URL")
	viper.BindEnv("messaging.rabbitmq.exchange", "RADIOKING_RABBITMQ_EXCHANGE")
	viper.BindEnv("messaging.rabbitmq.queue", "RADIOKING_RABBITMQ_QUEUE")
	viper.BindEnv("messaging.rabbitmq.routing_key", "RADIOKING_RABBITMQ_ROUTING_KEY")

	// Valeurs par défaut
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("auth.enabled", true)
	viper.SetDefault("auth.keycloak_url", "http://localhost:8180")
	viper.SetDefault("auth.realm", "radioking")
	viper.SetDefault("messaging.rabbitmq.url", "amqp://localhost:5672")
	viper.SetDefault("messaging.rabbitmq.exchange", "playlist_events")
	viper.SetDefault("messaging.rabbitmq.queue", "track_played")
	viper.SetDefault("messaging.rabbitmq.routing_key", "track.played")

	// Lecture optionnelle du fichier (ne fail pas si absent)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
