package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	HTTP     HttpConfig
	Database DatabaseConfig
}

// HttpConfig holds the HTTP server settings
type HttpConfig struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
}

// DatabaseConfig holds the database connection settings
type DatabaseConfig struct {
	URL               string        `env:"DATABASE_URL"`
	MaxOpenConns      int           `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	MaxIdleConns      int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime   time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	MaxConnIdleTime   time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`
	HealthCheckPeriod time.Duration `env:"DB_HEALTH_CHECK_PERIOD" envDefault:"1m"`
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

	// Set defaults for HTTP server
	vi.SetDefault("PORT", "8080")
	vi.SetDefault("HTTP_READ_TIMEOUT", "10s")
	vi.SetDefault("HTTP_WRITE_TIMEOUT", "10s")
	vi.SetDefault("HTTP_IDLE_TIMEOUT", "60s")

	// Set defaults for database connection pool
	vi.SetDefault("DB_MAX_OPEN_CONNS", 20)
	vi.SetDefault("DB_MAX_IDLE_CONNS", 5)
	vi.SetDefault("DB_CONN_MAX_LIFETIME", "5m")
	vi.SetDefault("DB_CONN_MAX_IDLE_TIME", "5m")
	vi.SetDefault("DB_HEALTH_CHECK_PERIOD", "1m")

	// Load database configuration from DATABASE_URL
	databaseURL := vi.GetString("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	dbConfig := DatabaseConfig{
		URL:               databaseURL,
		MaxOpenConns:      vi.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:      vi.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime:   vi.GetDuration("DB_CONN_MAX_LIFETIME"),
		MaxConnIdleTime:   vi.GetDuration("DB_CONN_MAX_IDLE_TIME"),
		HealthCheckPeriod: vi.GetDuration("DB_HEALTH_CHECK_PERIOD"),
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
