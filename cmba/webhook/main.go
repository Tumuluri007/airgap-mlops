// Package main implements the CMBA admission webhook.
//
// The webhook intercepts pod creation in the ml-serving namespace and
// performs six checks before admitting the pod (see handler.go).
//
// Air-gap mode: the webhook makes no external HTTPS calls. The Sigstore
// trust root and any required public keys are mounted as ConfigMaps on
// the webhook pod. Target latency: <50ms p99.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	listenAddr = flag.String("listen", ":8443", "TLS listen address for the webhook")
	tlsCert    = flag.String("tls-cert", "/etc/cmba/tls/tls.crt", "Path to TLS certificate")
	tlsKey     = flag.String("tls-key", "/etc/cmba/tls/tls.key", "Path to TLS key")
	trustRoot  = flag.String("trust-root", "/etc/cmba/trust-root", "Path to Sigstore offline trust root")
)

func main() {
	flag.Parse()

	h := NewHandler(*trustRoot)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", h.Validate)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	srv := &http.Server{
		Addr:              *listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	go func() {
		log.Printf("CMBA admission webhook listening on %s (trust-root=%s)", *listenAddr, *trustRoot)
		if err := srv.ListenAndServeTLS(*tlsCert, *tlsKey); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
