package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		systemCall := <-exit
		log.Printf("System call: %+v", systemCall)
		cancel()
	}()

	var config EnvConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to parse app config: %s", err)
	}
	if err := config.Validate(); err != nil {
		log.Fatalf("Failed to validate config: %s", err)
	}

	r := New(Config{Timeout: 10 * time.Second},
		config.AWSRegion,
		config.S3Bucket,
		config.S3Item,
	)

	srv := NewServer(config.Ports.HTTP, r)

	go func() {
		log.Print("Starting http server..")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen and serve: %s", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.Print("Stopping the http server..")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Could not shut down gracefully: %s", err)
		os.Exit(1)
	}

	defer cancel()
}

func NewServer(addr string, h http.Handler) *http.Server {
	s := &http.Server{Addr: addr, Handler: h}

	return s
}
