package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SEPP struct {
		FQDN string `yaml:"FQDN"`
		Port string `yaml:"Port"`
		TLS  struct {
			Enabled bool   `yaml:"Enabled"`
			Cert    string `yaml:"Cert"`
			Key     string `yaml:"Key"`
		} `yaml:"TLS"`
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

	if config.SEPP.TLS.Enabled && config.SEPP.TLS.Cert == "" && config.SEPP.TLS.Key == "" {
		return nil, fmt.Errorf("TLS.Cert and TLS.Key are required")
	}

	return &config, nil
}
