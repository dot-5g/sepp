package n32

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(certPath string, keyPath string, caCertPath string) *Client {
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

func (c *Client) POSTExchangeCapability(ctx context.Context, remoteURL string, secNegotiateReqData SecNegotiateReqData) (SecNegotiateRspData, error) {
	secNegotiateRspData := SecNegotiateRspData{}
	jsonData, err := json.Marshal(secNegotiateReqData)
	if err != nil {
		return secNegotiateRspData, err
	}

	endpoint := remoteURL + "/n32c-handshake/v1/exchange-capability"
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return secNegotiateRspData, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return secNegotiateRspData, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return secNegotiateRspData, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&secNegotiateRspData)
	if err != nil {
		return secNegotiateRspData, err
	}

	return secNegotiateRspData, nil
}
