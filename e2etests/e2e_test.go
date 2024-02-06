package e2etests

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"e2etests/certificates"
	"e2etests/docker"
)

const PLMNASBIFQDN = "https://0.0.0.0:1232"
const PLMNASEPPHostname = "sepp-plmn-a"
const PLMNBSEPPHostname = "sepp-plmn-b"
const DockerNetworkName = "n32"
const DockerImageName = "sepp-rock:test"

func Setup(projectPath string, seppAHostname string, seppBHostname string, dockerNetworkName string) {
	var err error
	PLMNACertsPath := projectPath + "/plmnA/certs/"
	PLMNAConfigPath := projectPath + "/plmnA/config.yaml"
	PLMNBCertsPath := projectPath + "/plmnB/certs/"
	PLMNBConfigPath := projectPath + "/plmnB/config.yaml"
	ClientCertsPath := projectPath + "/client/certs/"
	certificates.GenerateCertificates(PLMNACertsPath, seppAHostname, PLMNBCertsPath, seppBHostname, ClientCertsPath)
	if err := docker.CreateNetwork(dockerNetworkName); err != nil {
		log.Fatalf("Failed to create Docker network: %v", err)
	}
	if err = docker.RunContainer(seppAHostname, dockerNetworkName, DockerImageName, PLMNAConfigPath, PLMNACertsPath, map[string]string{"1231": "1231", "1232": "1232"}); err != nil {
		log.Fatalf("Failed to run PLMN A container: %v", err)
	}
	if err = docker.RunContainer(seppBHostname, dockerNetworkName, DockerImageName, PLMNBConfigPath, PLMNBCertsPath, map[string]string{"1233": "1233", "1234": "1234"}); err != nil {
		log.Fatalf("Failed to run PLMN B container: %v", err)
	}
}

func Cleanup(seppAID string, seppBID string, networkName string) {
	err := docker.StopAndRemoveContainer(seppAID)
	if err != nil {
		log.Printf("Failed to stop and remove container %s: %v", seppAID, err)
	}
	err = docker.StopAndRemoveContainer(seppBID)
	if err != nil {
		log.Printf("Failed to stop and remove container %s: %v", seppBID, err)
	}
	err = docker.RemoveNetwork(networkName)
	if err != nil {
		log.Printf("Failed to remove network %s: %v", networkName, err)
	}
}

func waitForService(client *http.Client, url string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(url)
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

func getClient(caCertPath string, clientCertPath string, clientKeyPath string) (*http.Client, error) {
	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load client certificate: %v", err)
	}
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("Failed to append CA certificate")
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
	return client, nil
}

func TestEndToEnd(t *testing.T) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	Setup(projectPath, PLMNASEPPHostname, PLMNBSEPPHostname, DockerNetworkName)
	defer Cleanup(PLMNASEPPHostname, PLMNBSEPPHostname, DockerNetworkName)
	client, err := getClient(projectPath+"/client/certs/ca.crt", projectPath+"/client/certs/client.crt", projectPath+"/client/certs/client.key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if err := waitForService(client, PLMNASBIFQDN, 10); err != nil {
		t.Fatalf("Failed to connect to SEPP in PLMN A: %v", err)
	}
}
