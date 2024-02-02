package n32

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/dot-5g/sepp/internal/model"
)

type N32C struct {
	FQDN         model.FQDN
	Capabilities []model.SecurityCapability
}

func loadClientCAs(caCertPath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request received: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("request handled: %s %s", r.Method, r.URL.Path)
	}
}

func StartServer(address string, serverCertPath string, serverKeyPath string, caCertPath string, fqdn string) {
	n32c := N32C{FQDN: model.FQDN(fqdn), Capabilities: []model.SecurityCapability{model.TLS}}
	http.HandleFunc("/n32c-handshake/v1/exchange-capability", loggingMiddleware(n32c.HandlePostExchangeCapability))
	clientCAPool, err := loadClientCAs(caCertPath)
	if err != nil {
		log.Fatalf("failed to load client CA certificate: %s", err)
	}
	tlsConfig := &tls.Config{
		ClientCAs:  clientCAPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
	}
	log.Printf("starting N32 server on %s", address)
	if err := server.ListenAndServeTLS(serverCertPath, serverKeyPath); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err)
	}
	log.Println("N32 server stopped")
}
