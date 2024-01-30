package n32

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(certPath, keyPath, caCertPath string) *Client {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}
}

func (c *Client) POSTExchangeCapability(remoteURL string, cap []SecurityCapability) (SecNegotiateRspData, error) {
	secNegotiateRspData := SecNegotiateRspData{}

	reqData := SecNegotiateReqData{
		// Populate reqData fields
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return secNegotiateRspData, err
	}

	req, err := http.NewRequest("POST", remoteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return secNegotiateRspData, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return secNegotiateRspData, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&secNegotiateRspData)
	if err != nil {
		return secNegotiateRspData, err
	}

	return secNegotiateRspData, nil
}
