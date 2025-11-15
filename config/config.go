package config

import (
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
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"postgres"`
	Password        string        `env:"DB_PASSWORD" envDefault:"password"`
	Database        string        `env:"DB_NAME" envDefault:"postgres"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
}

// NewConfig returns a new instance of Config
func NewConfig() *Config {
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

	// Set defaults for database
	vi.SetDefault("DB_HOST", "localhost")
	vi.SetDefault("DB_PORT", 5432)
	vi.SetDefault("DB_USER", "postgres")
	vi.SetDefault("DB_PASSWORD", "password")
	vi.SetDefault("DB_NAME", "postgres")
	vi.SetDefault("DB_SSLMODE", "disable")
	vi.SetDefault("DB_MAX_OPEN_CONNS", 25)
	vi.SetDefault("DB_MAX_IDLE_CONNS", 25)
	vi.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	// Load environment variables into struct
	return &Config{
		HTTP: HttpConfig{
			Port:         vi.GetString("PORT"),
			ReadTimeout:  vi.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: vi.GetDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  vi.GetDuration("HTTP_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            vi.GetString("DB_HOST"),
			Port:            vi.GetInt("DB_PORT"),
			User:            vi.GetString("DB_USER"),
			Password:        vi.GetString("DB_PASSWORD"),
			Database:        vi.GetString("DB_NAME"),
			SSLMode:         vi.GetString("DB_SSLMODE"),
			MaxOpenConns:    vi.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    vi.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: vi.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
	}
}
