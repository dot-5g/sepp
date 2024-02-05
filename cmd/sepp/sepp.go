package main

import (
	"flag"
	"log"
	"sync"

	"github.com/dot-5g/sepp/config"
	"github.com/dot-5g/sepp/internal/model"
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
		log.Fatalf("failed to read config file: %s", err)
	}
	seppContext := &model.SEPPContext{
		Mu:                 sync.Mutex{},
		LocalN32FQDN:       model.FQDN(conf.SEPP.Local.N32.FQDN),
		RemoteN32FQDN:      model.FQDN(""), // Initially empty
		SecurityCapability: model.SecurityCapability("TLS"),
	}

	startN32Server(&wg, conf.SEPP.Local.N32, seppContext)
	startSBIServer(&wg, conf.SEPP.Local.SBI, seppContext) // Always start SBI server

	// Capability exchange and updating SEPP context happens here
	remoteURL := conf.SEPP.Remote.URL
	if remoteURL != "" {
		secNegotiateRspData, err := exchangeCapability(remoteURL, conf.SEPP.Local.N32.FQDN, conf.SEPP.SecurityCapability, conf.SEPP.Remote.TLS)
		if err != nil {
			log.Fatalf("failed to exchange capability: %s", err)
		}
		seppContext.Mu.Lock()
		seppContext.RemoteN32FQDN = model.FQDN(secNegotiateRspData.Sender)
		seppContext.Mu.Unlock()
	}

	wg.Wait()
}

func startN32Server(wg *sync.WaitGroup, n32Config config.N32, seppContext *model.SEPPContext) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		n32.StartServer(n32Config.GetAddress(), n32Config.TLS.Cert, n32Config.TLS.Key, n32Config.TLS.CA, n32Config.FQDN, seppContext)
	}()
}

func startSBIServer(wg *sync.WaitGroup, sbiConfig config.SBI, seppContext *model.SEPPContext) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sbi.StartServer(sbiConfig.GetAddress(), sbiConfig.TLS.Cert, sbiConfig.TLS.Key, sbiConfig.TLS.CA, seppContext)
	}()
}

func exchangeCapability(remoteURL string, fqdn string, securityCapability string, n32TLSConf config.TLS) (n32.SecNegotiateRspData, error) {
	seppClient := n32.NewClient(n32TLSConf.Cert, n32TLSConf.Key, n32TLSConf.CA)
	reqData := n32.SecNegotiateReqData{
		Sender:                     model.FQDN(fqdn),
		SupportedSecCapabilityList: []model.SecurityCapability{model.SecurityCapability(securityCapability)},
	}
	secNegotiateRspData, err := seppClient.POSTExchangeCapability(remoteURL, reqData)
	if err != nil {
		log.Printf("failed to exchange capability: %s", err)
		return secNegotiateRspData, err
	}
	if secNegotiateRspData.SelectedSecCapability != model.TLS {
		log.Printf("failed to exchange capability: expected %s, got %s", model.TLS, secNegotiateRspData)
		return secNegotiateRspData, err
	}
	return secNegotiateRspData, nil
}
