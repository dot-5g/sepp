package n32

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/dot-5g/sepp/config"
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

func StartServer(config *config.Config) {
	n32c := N32C{FQDN: FQDN(config.SEPP.Local.N32.FQDN)}
	http.HandleFunc("/n32c-handshake/v1/exchange-capability", loggingMiddleware(n32c.HandlePostExchangeCapability))
	address := config.SEPP.Local.N32.Host + ":" + config.SEPP.Local.N32.Port
	clientCAPool, err := loadClientCAs(config.SEPP.Local.N32.TLS.CA)
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
	if err := server.ListenAndServeTLS(config.SEPP.Local.N32.TLS.Cert, config.SEPP.Local.N32.TLS.Key); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err)
	}
}
