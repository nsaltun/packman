package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP HttpConfig
	//TODO DatabaseConfig
}

type HttpConfig struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
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

	// Load environment variables into struct
	return &Config{
		HTTP: HttpConfig{
			Port:         vi.GetString("PORT"),
			ReadTimeout:  vi.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: vi.GetDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  vi.GetDuration("HTTP_IDLE_TIMEOUT"),
		},
	}
}
