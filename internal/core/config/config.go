package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TimeZone *time.Location `envconfig:"TIME_ZONE" default:"UTC"`
}

func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("CORE", &config); err != nil {
		return nil, fmt.Errorf("failed to process env core config: %w", err)
	}

	return &config, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()

	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
