package config

import (
	"github.com/Netflix/go-env"

	"github.com/mazxaxz/hls-streamer/pkg/logger"
	"github.com/mazxaxz/hls-streamer/pkg/rest"
)

type Config struct {
	HTTP   rest.Config   `env:"HTTP,required=true"`
	Logger logger.Config `env:"LOGGER"`
}

func Load() (Config, error) {
	var cfg Config
	if _, err := env.UnmarshalFromEnviron(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
