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
	conf, err := loadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	go func() {
		sbi.StartServer(conf)
	}()

	n32.StartServer(conf)
}

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func loadConfiguration(filePath string) (*config.Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf, err := config.ReadConfig(file)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
