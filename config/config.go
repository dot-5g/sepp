package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SEPP struct {
		Local struct {
			N32 struct {
				FQDN string `yaml:"fqdn"`
				Host string `yaml:"host"`
				Port string `yaml:"port"`
				TLS  struct {
					Cert string `yaml:"cert"`
					Key  string `yaml:"key"`
					CA   string `yaml:"ca"`
				} `yaml:"tls"`
			} `yaml:"n32"`
			SBI struct {
				Host string `yaml:"host"`
				Port string `yaml:"port"`
				TLS  struct {
					Cert string `yaml:"cert"`
					Key  string `yaml:"key"`
					CA   string `yaml:"ca"`
				} `yaml:"tls"`
			} `yaml:"sbi"`
		} `yaml:"local"`
		Remote struct {
			URL string `yaml:"url"`
			TLS struct {
				Cert string `yaml:"cert"`
				Key  string `yaml:"key"`
				CA   string `yaml:"ca"`
			} `yaml:"tls"`
		} `yaml:"remote"`
	} `yaml:"sepp"`
}

func ReadConfig(reader io.Reader) (*Config, error) {
	var config Config

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config data: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.SEPP.Local.N32.FQDN == "" {
		return fmt.Errorf("missing FQDN")
	}

	if config.SEPP.Local.N32.Host == "" {
		return fmt.Errorf("missing host")
	}

	if config.SEPP.Local.N32.Port == "" {
		return fmt.Errorf("missing port")
	}

	if config.SEPP.Local.N32.TLS.Cert == "" {
		return fmt.Errorf("missing TLS cert")
	}

	if config.SEPP.Local.N32.TLS.Key == "" {
		return fmt.Errorf("missing TLS key")
	}

	if config.SEPP.Local.N32.TLS.CA == "" {
		return fmt.Errorf("missing TLS CA")
	}

	if config.SEPP.Remote.URL != "" {
		if config.SEPP.Remote.TLS.Cert == "" {
			return fmt.Errorf("missing remote TLS cert")
		}

		if config.SEPP.Remote.TLS.Key == "" {
			return fmt.Errorf("missing remote TLS key")
		}

		if config.SEPP.Remote.TLS.CA == "" {
			return fmt.Errorf("missing remote TLS CA")
		}
	}
	return nil
}

func LoadConfiguration(filePath string) (*Config, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	conf, err := ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
