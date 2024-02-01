package sbi

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dot-5g/sepp/config"
)

func newReverseProxy(targetURL string) *httputil.ReverseProxy {
	url, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	return httputil.NewSingleHostReverseProxy(url)
}

func StartServer(config *config.Config) {
	proxy := newReverseProxy(config.SEPP.Remote.URL)
	http.Handle("/", proxy)

	address := config.SEPP.Local.SBI.Host + ":" + config.SEPP.Local.SBI.Port
	server := &http.Server{
		Addr:    address,
		Handler: proxy,
	}

	log.Printf("Starting SBI server on %s", address)
	if err := server.ListenAndServeTLS(config.SEPP.Local.SBI.TLS.Cert, config.SEPP.Local.SBI.TLS.Key); err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("SBI server stopped")
}
