package app

import "glide/internal"

const (
	LoadFileUrl  = "media/"
	DefaultImage = ""
)

type Microservice struct {
	SessionServerUrl string `toml:"session_url"`
}

type RepositoryConnections struct {
	DataBaseUrl     string `toml:"database_url"`
	SessionRedisUrl string `toml:"session-redis_url"`
	RabbitUrl       string `toml:"rabbit_url"`
}

type Config struct {
	internal.Config
	MediaDir      string                `toml:"media_dir"`
	Microservices Microservice          `toml:"microservice"`
	Repository    RepositoryConnections `toml:"connections"`
	Cors          internal.CorsConfig   `toml:"cors"`
}

func NewConfig() *Config {
	return &Config{}
}
