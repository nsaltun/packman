package config

import "time"

type Config struct {
	HTTP HttpConfig
	//TODO DatabaseConfig
}

type HttpConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewConfig() *Config {
	return &Config{
		HTTP: HttpConfig{
			Port:         "8080",
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			IdleTimeout:  time.Second * 60,
		},
	}
}
