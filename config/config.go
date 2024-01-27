package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SEPP struct {
		FQDN string `yaml:"FQDN"`
		Host string `yaml:"Host"`
		Port string `yaml:"Port"`
		TLS  struct {
			Cert string `yaml:"Cert"`
			Key  string `yaml:"Key"`
			CA   string `yaml:"CA"`
		} `yaml:"TLS"`
	} `yaml:"SEPP"`
	NRF struct {
		FQDN string `yaml:"FQDN"`
		Port string `yaml:"Port"`
	}
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

	if config.SEPP.TLS.Cert == "" {
		return nil, fmt.Errorf("SEPP.TLS.Cert is required")
	}

	if config.SEPP.TLS.Key == "" {
		return nil, fmt.Errorf("SEPP.TLS.Key is required")
	}

	if config.SEPP.TLS.CA == "" {
		return nil, fmt.Errorf("SEPP.TLS.CACert is required")
	}

	return &config, nil
}
