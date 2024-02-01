package nsepp

import (
	"log"
	"net/http"

	"github.com/dot-5g/sepp/config"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("sbi - request received: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("sbi - request handled: %s %s", r.Method, r.URL.Path)
	}
}

func StartServer(address string, config *config.TLS) {
	server := &http.Server{
		Addr: address,
	}
	http.HandleFunc("/nsepp-telescopic/v1/mapping", loggingMiddleware(HandleGetMapping))
	log.Printf("starting SBI server on %s", address)
	if err := server.ListenAndServeTLS(config.Cert, config.Key); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println("SBI server stopped")
}
