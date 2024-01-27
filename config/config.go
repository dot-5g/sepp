package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SEPP struct {
		FQDN string `yaml:"FQDN"`
	} `yaml:"SEPP"`
}

func ReadConfig(configPath string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.SEPP.FQDN == "" {
		return nil, fmt.Errorf("FQDN is required")
	}

	return &config, nil
}