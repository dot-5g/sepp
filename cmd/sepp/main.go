package main

import (
	"flag"
	"os"

	"github.com/dot-5g/sepp/config"

	"github.com/dot-5g/sepp/internal/n32"
	"github.com/dot-5g/sepp/internal/sbi"

	"log"
)

var configFilePath string

func main() {
	flag.Parse()
	config, err := loadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	go func() {
		sbi.StartServer(config)
	}()

	n32.StartServer(config)
}

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func loadConfiguration(filePath string) (*config.Config, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	conf, err := config.ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
