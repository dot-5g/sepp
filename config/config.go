package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type TLS struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	CA   string `yaml:"ca"`
}

type N32 struct {
	FQDN string `yaml:"fqdn"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	TLS  TLS    `yaml:"tls"`
}

type SBI struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	TLS  TLS    `yaml:"tls"`
}

type Local struct {
	N32 N32 `yaml:"n32"`
	SBI SBI `yaml:"sbi"`
}

type Remote struct {
	URL string `yaml:"url"`
	TLS TLS    `yaml:"tls"`
}

type SEPP struct {
	SecurityCapability string `yaml:"securityCapability"`
	Local              Local  `yaml:"local"`
	Remote             Remote `yaml:"remote"`
}

type Config struct {
	SEPP SEPP `yaml:"sepp"`
}

func (n32 N32) GetAddress() string {
	return n32.Host + ":" + n32.Port
}

func (sbi SBI) GetAddress() string {
	return sbi.Host + ":" + sbi.Port
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

	if config.SEPP.SecurityCapability != "TLS" {
		return fmt.Errorf("unsupported security capability, only TLS is supported")
	}

	if config.SEPP.Local.N32.FQDN == "" {
		return fmt.Errorf("missing Local N32 FQDN")
	}

	if config.SEPP.Local.N32.Host == "" {
		return fmt.Errorf("missing Local N32 Host")
	}

	if config.SEPP.Local.N32.Port == "" {
		return fmt.Errorf("missing port")
	}

	if config.SEPP.Local.N32.TLS.Cert == "" {
		return fmt.Errorf("missing Local N32 TLS Cert")
	}

	if config.SEPP.Local.N32.TLS.Key == "" {
		return fmt.Errorf("missing Local N32 TLS Key")
	}

	if config.SEPP.Local.N32.TLS.CA == "" {
		return fmt.Errorf("missing Local N32 TLS CA")
	}

	if config.SEPP.Remote.URL != "" {
		if config.SEPP.Remote.TLS.Cert == "" {
			return fmt.Errorf("missing Remote TLS Cert")
		}

		if config.SEPP.Remote.TLS.Key == "" {
			return fmt.Errorf("missing Remote TLS Key")
		}

		if config.SEPP.Remote.TLS.CA == "" {
			return fmt.Errorf("missing Remote TLS CA")
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
