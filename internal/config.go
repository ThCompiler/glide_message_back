package internal

type Config struct {
	LogLevel string `toml:"log_level"`
	LogAddr  string `toml:"log_path"`
	Port     string `toml:"port"`
}

type CorsConfig struct {
	Urls    []string `toml:"urls"`
	Headers []string `toml:"headers"`
	Methods []string `toml:"methods"`
}
