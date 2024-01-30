package main

import (
	"flag"
	"os"
	"time"

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

	address := config.SEPP.Local.N32.Host + ":" + config.SEPP.Local.N32.Port
	caCertPath := config.SEPP.Local.N32.TLS.CA
	serverCertPath := config.SEPP.Local.N32.TLS.Cert
	serverKeyPath := config.SEPP.Local.N32.TLS.Key
	fqdn := config.SEPP.Local.N32.FQDN

	go n32.StartServer(address, serverCertPath, serverKeyPath, caCertPath, fqdn)

	seppClient := n32.NewClient(config.SEPP.Local.N32.TLS.Cert, config.SEPP.Local.N32.TLS.Key, config.SEPP.Local.N32.TLS.CA)

	remoteURL := config.SEPP.Remote.URL
	securityCapabilities := []n32.SecurityCapability{n32.TLS}

	for {
		cap, err := seppClient.POSTExchangeCapability(remoteURL, securityCapabilities)
		if err != nil {
			log.Printf("failed to exchange capability: %s", err)
			time.Sleep(30 * time.Second) // Retry after some time
			continue
		}

		if cap.SelectedSecCapability == n32.TLS {
			log.Println("security exchange successful, starting SBI server...")
			sbi.StartServer(config)
			break
		} else {
			log.Printf("unsupported capability: %v", cap)
			time.Sleep(30 * time.Second) // Retry after some time
		}
	}

	// Keep the main goroutine running
	select {}
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
