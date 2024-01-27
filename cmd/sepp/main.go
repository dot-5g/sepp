package main

import (
	"flag"

	"github.com/dot-5g/sepp/config"

	"github.com/dot-5g/sepp/internal/server"

	"log"
)

var configFilePath string

func main() {
	flag.Parse()
	config, err := loadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	server.Start(config)
}

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func loadConfiguration(filePath string) (*config.Config, error) {
	conf, err := config.ReadConfig(filePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
