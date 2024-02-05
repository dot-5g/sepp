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
func dynamicProxyHandler(seppContext *model.SEPPContext) http.HandlerFunc {
	var mu sync.Mutex // Protects access to the reverseProxy to ensure thread safety
	var reverseProxy *httputil.ReverseProxy

	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if r.TLS == nil {
			http.Error(w, "TLS is required", http.StatusBadRequest)
			log.Println("SBI server - TLS is required")
			return
		}

		if r.TLS.PeerCertificates == nil {
			http.Error(w, "Client certificate is required", http.StatusBadRequest)
			log.Println("SBI server - client certificate is required")
			return
		}

		remoteURL := string(seppContext.RemoteN32FQDN)
		if remoteURL == "" {
			http.Error(w, "Remote SEPP not configured", http.StatusInternalServerError)
			log.Println("SBI server - remote SEPP not configured")
			return
		}

		if reverseProxy == nil || reverseProxy.Director == nil || r.URL.Host != remoteURL {
			targetURL, err := url.Parse(remoteURL)
			if err != nil {
				http.Error(w, "Failed to parse target URL", http.StatusInternalServerError)
				log.Printf("SBI server - failed to parse target URL: %v", err)
				return
			}
			reverseProxy = httputil.NewSingleHostReverseProxy(targetURL)
			log.Printf("SBI server - forwarding requests to remote SEPP (%s)", remoteURL)
		}

		reverseProxy.ServeHTTP(w, r)
	}
}

func StartServer(address string, serverCertPath string, serverKeyPath string, caCertPath string, seppContext *model.SEPPContext) {
	cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		log.Fatalf("failed to load server key pair: %v", err)
	}

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", dynamicProxyHandler(seppContext))

	log.Printf("SBI server - started listening on %s", address)
	if err := server.ListenAndServeTLS(serverCertPath, serverKeyPath); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println("SBI server - stopped")
}
