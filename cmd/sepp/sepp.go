package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dot-5g/sepp/config"

	"github.com/dot-5g/sepp/internal/n32"
	"github.com/dot-5g/sepp/internal/sbi"

	"log"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func main() {
	flag.Parse()

	conf, err := config.LoadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	address := conf.SEPP.Local.N32.Host + ":" + conf.SEPP.Local.N32.Port
	go n32.StartServer(address, conf.SEPP.Local.N32.TLS.Cert, conf.SEPP.Local.N32.TLS.Key, conf.SEPP.Local.N32.TLS.CA, conf.SEPP.Local.N32.FQDN)
	remoteURL := conf.SEPP.Remote.URL

	if remoteURL == "" {
		log.Println("no remote URL specified, not starting SBI server...")
		waitForShutdown()
	}

	seppClient := n32.NewClient(conf.SEPP.Local.N32.TLS.Cert, conf.SEPP.Local.N32.TLS.Key, conf.SEPP.Local.N32.TLS.CA)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
		log.Println("server gracefully stopped")
	}()

	reqData := n32.SecNegotiateReqData{
		Sender:                     n32.FQDN("testSender"),
		SupportedSecCapabilityList: []n32.SecurityCapability{n32.TLS},
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if cap, err := seppClient.POSTExchangeCapability(remoteURL, reqData); err != nil {
				log.Printf("failed to exchange capability: %s", err)
				waitOrCancel(ctx, 30*time.Second)
			} else if cap.SelectedSecCapability == n32.TLS {
				log.Println("security exchange successful, starting SBI server...")
				sbi.StartServer(conf)
				return
			} else {
				log.Printf("unsupported capability: %v", cap)
				waitOrCancel(ctx, 30*time.Second)
			}
		}
	}
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("server gracefully stopped")
}

func waitOrCancel(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(duration):
		return
	}
}
