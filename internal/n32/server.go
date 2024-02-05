package n32

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/dot-5g/sepp/internal/model"
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
		log.Printf("N32 server - request received: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("N32 server - request handled: %s %s", r.Method, r.URL.Path)
	}
}

func HandleN32f(w http.ResponseWriter, r *http.Request) {
	log.Printf("N32 server - request received: %s %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	log.Printf("N32 server - request handled: %s %s", r.Method, r.URL.Path)
}

func StartServer(address string, serverCertPath string, serverKeyPath string, caCertPath string, fqdn string, seppContext *model.SEPPContext) {
	mux := http.NewServeMux()
	mux.HandleFunc("/n32c-handshake/v1/exchange-capability", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		HandlePostExchangeCapability(w, r, seppContext)
	}))
	mux.HandleFunc("/", HandleN32f)
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
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
	log.Printf("N32 server - started listening on %s", address)
	if err := server.ListenAndServeTLS(serverCertPath, serverKeyPath); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %s", err)
	}
	log.Println("N32 server - stopped")
}
