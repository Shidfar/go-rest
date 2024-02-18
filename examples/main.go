package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shidfar/go-rest/examples/fooservice"
	"github.com/Shidfar/go-rest/examples/fooservice/fooservicehttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a new instance of a string service.
	svc := fooservice.New()

	log.Println("starting string service on localhost:8080")

	// Register http handlers for the service.
	// In this case,
	// one http handler will be registered [GET/POST] /count.
	// Get usage: /count?text=hello
	// Post usage: /count -d '{"text":"hello"}'
	fooservicehttp.RegisterHandlers(http.DefaultServeMux, svc, "/api")

	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	if err := waitForShutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func waitForShutdown(ctx context.Context) error {
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to receive the signal when the context expires
	done := make(chan struct{})

	// Run a goroutine that closes the done channel when the context expires
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()

	// Wait for a signal or context expiration
	select {
	case <-sigs:
		fmt.Println("Received an interrupt signal, shutting down...")
		return nil
	case <-done:
		fmt.Println("Context has expired, shutting down...")
		return errors.New("context expired")
	}
}
