package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	portHTTP := "8080"
	portHTTPS := "8443"
	certFile := getenv("SSL_CERT_FILE", "")
	keyFile := getenv("SSL_KEY_FILE", "")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	})

	// If both cert and key are set, serve ONLY HTTPS.
	if certFile != "" && keyFile != "" {
		addr := ":" + portHTTPS
		log.Printf("Starting HTTPS ONLY on %s (cert=%s key=%s)", addr, certFile, keyFile)
		if err := http.ListenAndServeTLS(addr, certFile, keyFile, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTPS server error: %v", err)
		}
		return
	}

	// Otherwise, serve ONLY HTTP.
	addr := ":" + portHTTP
	log.Printf("TLS envs not set -> starting HTTP ONLY on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
