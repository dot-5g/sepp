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
		log.Fatalf("failed to parse target URL: %v", err)
	}

	return httputil.NewSingleHostReverseProxy(url)
}

func StartServer(remoteURL string, address string, sbiTLS *config.TLS) {
	proxy := newReverseProxy(remoteURL)
	http.Handle("/", proxy)

	server := &http.Server{
		Addr:    address,
		Handler: proxy,
	}
	log.Printf("starting SBI server on %s", address)
	if err := server.ListenAndServeTLS(sbiTLS.Cert, sbiTLS.Key); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println("SBI server stopped")
}
