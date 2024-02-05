package main_test

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"testing"
)

var PLMNASBIFQDN = "https://0.0.0.0:1232"
var CertsPath = "client/certs/"

func TestEndToEnd(t *testing.T) {
	clientCert, err := tls.LoadX509KeyPair(CertsPath+"client.crt", CertsPath+"client.key")
	if err != nil {
		t.Fatalf("Failed to load client certificate: %v", err)
	}
	caCert, err := os.ReadFile(CertsPath + "ca.crt")
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
