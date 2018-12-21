package app

import (
	"fmt"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type Config struct {
	CACertPath string `env:"CA_CERT_PATH"`
	CertPath   string `env:"CERT_PATH"`
	KeyPath    string `env:"KEY_PATH"`
	ShardID    string `env:"SHARD_ID"`
	LoggrAddr  string `env:"LOGGR_ADDR"`
	RabbitAddr string `env:"RABBIT_ADDR"`
}

// LoadConfig reads from the environment to create a Config.
func LoadConfig() (*Config, error) {
	config := Config{}
	err := envstruct.Load(&config)
	if err != nil {
		return nil, err
	}

	if config.LoggrAddr == "" {
		return nil, fmt.Errorf("Loggregator address is required")
	}

	if config.RabbitAddr == "" {
		return nil, fmt.Errorf("RabbitMQ address is required")
	}

	return &config, nil
}
