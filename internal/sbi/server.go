package sbi

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	"github.com/dot-5g/sepp/internal/model"
)

// dynamicProxyHandler creates a handler function that dynamically decides
// the target URL based on the current state of seppContext.RemoteN32FQDN.
func dynamicProxyHandler(seppContext *model.SEPPContext, outboundTLSConfig *tls.Config) http.HandlerFunc {
	var mu sync.Mutex
	var reverseProxy *httputil.ReverseProxy

	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		remoteURL := string(seppContext.RemoteN32FQDN)
		if remoteURL == "" {
			http.Error(w, "Remote SEPP not configured", http.StatusInternalServerError)
			return
		}

		if reverseProxy == nil {
			targetURL, err := url.Parse(remoteURL)
			if err != nil {
				http.Error(w, "Failed to parse target URL", http.StatusInternalServerError)
				return
			}
			reverseProxy = httputil.NewSingleHostReverseProxy(targetURL)
			reverseProxy.Transport = &http.Transport{
				TLSClientConfig: outboundTLSConfig,
			}
			log.Printf("SBI server - forwarding requests to remote SEPP (%s)", remoteURL)
		} else {
			log.Printf("SBI server - reusing existing reverse proxy to remote SEPP (%s)", remoteURL)
		}

		reverseProxy.ServeHTTP(w, r)
	}
}

func StartServer(address, serverCertPath, serverKeyPath, caCertPath, clientCertPath, clientKeyPath string, seppContext *model.SEPPContext) {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to append CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("failed to load client certificate and key: %v", err)
	}

	outboundTLSConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", dynamicProxyHandler(seppContext, outboundTLSConfig))

	serverCert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		log.Fatalf("failed to load server key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      address,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Printf("SBI server - started listening on %s", address)
	if err := server.ListenAndServeTLS(serverCertPath, serverKeyPath); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println("SBI server - stopped")
}
