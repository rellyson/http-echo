package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rellyson/http-echo/internal/handlers"
	"github.com/rellyson/http-echo/pkg/middlewares"
)

var (
	listenFlag          = flag.String("listen", ":3000", "Address to listen on")
	metricsPathFlag     = flag.String("metrics", "/metrics", "Path to expose metrics")
	defaultReadTimeout  = 3 * time.Second
	defaultWriteTimeout = 3 * time.Second
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("GET %s", *metricsPathFlag), handlers.Metrics())
	mux.HandleFunc("GET /health", handlers.HealthCheck)
	mux.HandleFunc("/", handlers.Echo)

	mw := middlewares.CreateStack(
		middlewares.Recover,
		middlewares.Logging,
	)
	server := &http.Server{
		Addr:         *listenFlag,
		Handler:      mw(mux),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	go func() {
		log.Printf("[INFO] server listening on %s", *listenFlag)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[ERROR] server exited with: %s", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt
	<-signalCh

	log.Printf("[INFO] received interrupt, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] failed to shutdown server: %s", err)
	}
}
