package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/dot-5g/sepp/config"
	"github.com/dot-5g/sepp/internal/n32"
	"github.com/dot-5g/sepp/internal/sbi"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	conf, err := config.LoadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		n32.StartServer(conf.SEPP.Local.N32.GetAddress(), conf.SEPP.Local.N32.TLS.Cert, conf.SEPP.Local.N32.TLS.Key, conf.SEPP.Local.N32.TLS.CA, conf.SEPP.Local.N32.FQDN)
	}()
	remoteURL := conf.SEPP.Remote.URL
	if remoteURL != "" {
		exchangeCapability(remoteURL, conf.SEPP.Remote.TLS)
		wg.Add(1)
		go func() {
			defer wg.Done()
			sbi.StartServer(conf)
		}()
	}
	wg.Wait()
}

func exchangeCapability(remoteURL string, n32TLSConf config.TLS) {
	seppClient := n32.NewClient(n32TLSConf.Cert, n32TLSConf.Key, n32TLSConf.CA)
	reqData := n32.SecNegotiateReqData{
		Sender:                     n32.FQDN("testSender"),
		SupportedSecCapabilityList: []n32.SecurityCapability{n32.TLS},
	}
	for {
		cap, err := seppClient.POSTExchangeCapability(remoteURL, reqData)
		if err == nil && cap.SelectedSecCapability == n32.TLS {
			log.Printf("Successfully exchanged capability: %s", cap.SelectedSecCapability)
			break
		}
		if err != nil {
			log.Printf("Failed to exchange capability: %s", err)
		} else {
			log.Printf("Failed to exchange capability: expected %s, got %s", n32.TLS, cap)
		}
		time.Sleep(30 * time.Second)
	}
}
