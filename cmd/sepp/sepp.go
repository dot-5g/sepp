package main

import (
	"flag"
	"log"
	"sync"

	"github.com/dot-5g/sepp/config"
	"github.com/dot-5g/sepp/internal/n32"
	"github.com/dot-5g/sepp/internal/nsepp"
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
		log.Fatalf("failed to read config file: %s", err)
	}
	startN32Server(&wg, conf.SEPP.Local.N32)
	remoteURL := conf.SEPP.Remote.URL
	if remoteURL != "" {
		exchangeCapability(remoteURL, conf.SEPP.Local.N32.FQDN, conf.SEPP.SecurityCapability, conf.SEPP.Remote.TLS)
	}
	startSBIServer(&wg, remoteURL, conf.SEPP.Local.SBI)
	startNSEPPServer(&wg, conf.SEPP.Local.NSEPP)
	wg.Wait()
}

func startN32Server(wg *sync.WaitGroup, n32Config config.N32) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		n32.StartServer(n32Config.GetAddress(), n32Config.TLS.Cert, n32Config.TLS.Key, n32Config.TLS.CA, n32Config.FQDN)
	}()
}

func startSBIServer(wg *sync.WaitGroup, remoteURL string, sbiConfig config.SBI) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sbi.StartServer(remoteURL, sbiConfig.GetAddress(), sbiConfig.TLS)
	}()
}

func startNSEPPServer(wg *sync.WaitGroup, nseppConfig config.NSEPP) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		nsepp.StartServer(nseppConfig.GetAddress(), nseppConfig.TLS)
	}()
}

func exchangeCapability(remoteURL string, fqdn string, securityCapability string, n32TLSConf config.TLS) {
	seppClient := n32.NewClient(n32TLSConf.Cert, n32TLSConf.Key, n32TLSConf.CA)
	reqData := n32.SecNegotiateReqData{
		Sender:                     n32.FQDN(fqdn),
		SupportedSecCapabilityList: []n32.SecurityCapability{n32.SecurityCapability(securityCapability)},
	}
	cap, err := seppClient.POSTExchangeCapability(remoteURL, reqData)
	if err == nil && cap.SelectedSecCapability == n32.TLS {
		log.Printf("successfully exchanged capability %s with remote SEPP %s", cap.SelectedSecCapability, remoteURL)
		return
	}
	if err != nil {
		log.Printf("failed to exchange capability: %s", err)
	} else {
		log.Printf("failed to exchange capability: expected %s, got %s", n32.TLS, cap)
	}
}
