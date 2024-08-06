package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	GRPC GRPC
	HTTP HTTP
	Log  Log
}

func NewConfig() (Config, error) {
	var (
		config = Config{}
		err    error
	)

	err = env.Parse(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse env: %w", err)
	}

	return config, nil
}
