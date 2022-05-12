package push

import "glide/internal"

type Config struct {
	internal.Config
	RabbitUrl  string              `toml:"rabbit_url"`
	Cors       internal.CorsConfig `toml:"cors"`
	SessionUrl string              `toml:"session_url"`
	SqlUrl     string              `toml:"database_url"`
}
