package e2e_tests

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"testing"
)

var PLMNAFQDN = "https://127.0.0.1:1234"
var CertsPath = "client/certs/"

func TestEndToEnd(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(CertsPath+"client.crt", CertsPath+"client.key")
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
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get(PLMNAFQDN)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

}