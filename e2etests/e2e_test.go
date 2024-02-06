package test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dot-5g/sepp/e2etests/certificates"
	"github.com/dot-5g/sepp/e2etests/docker"
)

const PLMNASBIFQDN = "https://0.0.0.0:1232"
const PLMNACertsPath = "plmnA/certs/"
const PLMNBCertsPath = "plmnB/certs/"
const ClientCertsPath = "client/certs/"
const PLMNASEPPHostname = "sepp-plmn-a"
const PLMNBSEPPHostname = "sepp-plmn-b"
const DockerNetworkName = "n32"

func Setup(seppAHostname string, seppBHostname string, dockerNetworkName string) {
	var err error
	certificates.GenerateCertificates(PLMNACertsPath, seppAHostname, PLMNBCertsPath, seppBHostname, ClientCertsPath)
	if err := docker.CreateNetwork(dockerNetworkName); err != nil {
		log.Fatalf("Failed to create Docker network: %v", err)
	}
	if err = docker.RunContainer(seppAHostname, dockerNetworkName, "sepp:0.1", "/home/guillaume/code/sepp/e2etests/plmnA/config.yaml", "/home/guillaume/code/sepp/e2etests/plmnA/certs/", map[string]string{"1231": "1231", "1232": "1232"}); err != nil {
		log.Fatalf("Failed to run PLMN A container: %v", err)
	}
	if err = docker.RunContainer(seppBHostname, dockerNetworkName, "sepp:0.1", "/home/guillaume/code/sepp/e2etests/plmnB/config.yaml", "/home/guillaume/code/sepp/e2etests/plmnB/certs/", map[string]string{"1233": "1233", "1234": "1234"}); err != nil {
		log.Fatalf("Failed to run PLMN B container: %v", err)
	}
}

func Cleanup(seppAID string, seppBID string, networkName string) {
	docker.StopAndRemoveContainer(seppAID)
	docker.StopAndRemoveContainer(seppBID)
	docker.RemoveNetwork(networkName)
}

func waitForService(url string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("Service available at %s", url)
				return nil
			}
		}
		log.Printf("Service not available at %s, retrying...", url)
		time.Sleep(time.Duration(i) * time.Second)
	}
	return fmt.Errorf("service not available at %s", url)
}

func TestEndToEnd(t *testing.T) {
	Setup(PLMNASEPPHostname, PLMNBSEPPHostname, DockerNetworkName)
	defer Cleanup(PLMNASEPPHostname, PLMNBSEPPHostname, DockerNetworkName)
	clientCert, err := tls.LoadX509KeyPair(ClientCertsPath+"client.crt", ClientCertsPath+"client.key")
	if err != nil {
		t.Fatalf("Failed to load client certificate: %v", err)
	}
	caCert, err := os.ReadFile(ClientCertsPath + "ca.crt")
	if err != nil {
		t.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		t.Fatal("Failed to append CA certificate")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	if err := waitForService(PLMNASBIFQDN, 10); err != nil {
		t.Fatalf("Failed to connect to SEPP in PLMN A: %v", err)
	}
	address := PLMNASBIFQDN
	resp, err := client.Get(address)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}
}
