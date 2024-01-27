package server

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/dot-5g/sepp/config"
	"github.com/dot-5g/sepp/internal/n32"
)

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

func Start(conf *config.Config) {
	n32c := n32.N32C{FQDN: n32.FQDN(conf.SEPP.FQDN)}
	http.HandleFunc("/n32c-handshake/v1/exchange-capability", loggingMiddleware(n32c.HandlePostExchangeCapability))
	address := conf.SEPP.Host + ":" + conf.SEPP.Port
	clientCAPool, err := loadClientCAs(conf.SEPP.TLS.CA)
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
	log.Printf("starting server on %s", address)
	if err := server.ListenAndServeTLS(conf.SEPP.TLS.Cert, conf.SEPP.TLS.Key); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err)
	}
}
