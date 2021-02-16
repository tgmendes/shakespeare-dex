package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tgmendes/shakespeare-dex/handler"
)

func main() {
	if err := run(); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run() error {
	// On Interrupt/SIGTERM tell net/http server to shutdown gracefully.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// these could be defined as environment variables, config values or command line params -
	// we hardcode them here for simplicity
	const (
		Host         = ":5000"
		ReadTimeout  = 5 * time.Second
		WriteTimeout = 5 * time.Second
	)

	server := http.Server{
		Addr:         Host,
		Handler:      handler.ShakespeareDexAPI(),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("%s: %v", "server error", err)

	case <-shutdown:
		log.Println("Starting API Shutdown")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", timeout, err)
			err = server.Close()
		}

		if err != nil {
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
	}

	return nil
}
