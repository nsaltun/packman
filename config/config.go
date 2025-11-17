package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HttpConfig
	Database DatabaseConfig
}

type HttpConfig struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
}

type DatabaseConfig struct {
	URL             string        `env:"DATABASE_URL"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
}

// NewConfig returns a new instance of Config
func NewConfig() (*Config, error) {
	vi := viper.New()

	// Set config file name and type
	vi.SetConfigName(".env")
	vi.SetConfigType("env")

	// Add paths to search for .env file
	vi.AddConfigPath(".")      // Current directory
	vi.AddConfigPath("./")     // Current directory
	vi.AddConfigPath("../")    // Parent directory
	vi.AddConfigPath("../../") // Two levels up

	// Read .env file if it exists (ignore error if not found)
	_ = vi.ReadInConfig()

	// AutomaticEnv makes viper check environment variables
	// Environment variables take precedence over .env file
	vi.AutomaticEnv()

	// Set defaults
	vi.SetDefault("PORT", "8080")
	vi.SetDefault("HTTP_READ_TIMEOUT", "10s")
	vi.SetDefault("HTTP_WRITE_TIMEOUT", "10s")
	vi.SetDefault("HTTP_IDLE_TIMEOUT", "60s")

	// Set defaults for database connection pool
	vi.SetDefault("DB_MAX_OPEN_CONNS", 25)
	vi.SetDefault("DB_MAX_IDLE_CONNS", 25)
	vi.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	// Load database configuration from DATABASE_URL
	databaseURL := vi.GetString("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	dbConfig := DatabaseConfig{
		URL:             databaseURL,
		MaxOpenConns:    vi.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    vi.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: vi.GetDuration("DB_CONN_MAX_LIFETIME"),
	}

	return &Config{
		HTTP: HttpConfig{
			Port:         vi.GetString("PORT"),
			ReadTimeout:  vi.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: vi.GetDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  vi.GetDuration("HTTP_IDLE_TIMEOUT"),
		},
		Database: dbConfig,
	}, nil
}
