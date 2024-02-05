package main

import (
	"flag"
	"log"
	"sync"
	"time"

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
		Mu:                          sync.Mutex{},
		LocalN32FQDN:                model.FQDN(conf.SEPP.Local.N32.FQDN),
		RemoteN32FQDN:               model.FQDN(""),
		SupportedSecurityCapability: model.SecurityCapability("TLS"),
	}
	startN32Server(&wg, conf.SEPP.Local.N32, seppContext)
	startSBIServer(&wg, conf.SEPP.Local.SBI, conf.SEPP.Remote.TLS, seppContext)
	exchangeCapability(conf.SEPP.Remote.URL, conf.SEPP.Local.N32.FQDN, conf.SEPP.SecurityCapability, conf.SEPP.Remote.TLS, seppContext)
	wg.Wait()
}

func startN32Server(wg *sync.WaitGroup, n32Config config.N32, seppContext *model.SEPPContext) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		n32.StartServer(n32Config.GetAddress(), n32Config.TLS.Cert, n32Config.TLS.Key, n32Config.TLS.CA, n32Config.FQDN, seppContext)
	}()
}

func startSBIServer(wg *sync.WaitGroup, sbiConfig config.SBI, clientTLS config.TLS, seppContext *model.SEPPContext) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sbi.StartServer(sbiConfig.GetAddress(), sbiConfig.TLS.Cert, sbiConfig.TLS.Key, sbiConfig.TLS.CA, clientTLS.Cert, clientTLS.Key, seppContext)
	}()
}

func exchangeCapability(remoteURL string, fqdn string, securityCapability string, n32TLSConf config.TLS, seppContext *model.SEPPContext) {
	for {
		seppContext.Mu.Lock()
		selectedCapability := seppContext.SelectedSecurityCapability
		seppContext.Mu.Unlock()
		if selectedCapability != "" {
			return
		}
		seppClient := n32.NewClient(n32TLSConf.Cert, n32TLSConf.Key, n32TLSConf.CA)
		reqData := n32.SecNegotiateReqData{
			Sender:                     model.FQDN(fqdn),
			SupportedSecCapabilityList: []model.SecurityCapability{model.SecurityCapability(securityCapability)},
		}
		secNegotiateRspData, err := seppClient.POSTExchangeCapability(remoteURL, reqData)
		if err != nil || secNegotiateRspData.SelectedSecCapability != model.TLS {
			time.Sleep(5 * time.Second)
			continue
		}
		seppContext.Mu.Lock()
		seppContext.RemoteN32FQDN = model.FQDN(secNegotiateRspData.Sender)
		seppContext.SelectedSecurityCapability = secNegotiateRspData.SelectedSecCapability
		seppContext.Mu.Unlock()
		return
	}
}
