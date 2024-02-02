package nsepp

import (
	"log"
	"net/http"

	"github.com/dot-5g/sepp/config"
	"github.com/dot-5g/sepp/internal/model"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("nsepp - request received: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("nsepp - request handled: %s %s", r.Method, r.URL.Path)
	}
}

func StartServer(address string, config config.TLS, seppContext *model.SEPPContext) {
	server := &http.Server{
		Addr: address,
	}

	http.HandleFunc("/nsepp-telescopic/v1/mapping", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		HandleGetMapping(w, r, seppContext)
	}))
	log.Printf("starting Nsepp server on %s", address)
	if err := server.ListenAndServeTLS(config.Cert, config.Key); err != http.ErrServerClosed {
		log.Fatalf("failed to start Nsepp server: %v", err)
	}
	log.Println("Nsepp server stopped")
}
