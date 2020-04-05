package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/internal/whatflixapis/version"
	"github.com/whatflix/pkg/httphealthcheck"
	"github.com/whatflix/pkg/utils"
)

const healthCheckAddr = "/movies/health-check"

func main() {
	err := run()
	if err != nil {
		handleError(err)
	}
}

func run() error {
	utils.InitLog()
	config, err := getConfigs()
	if err != nil {
		return errors.WithMessage(err, "get configs")
	}
	utils.LogStart(version.Version, config.Env)

	httpServer, err := initHTTPServer(config)
	if err != nil {
		return errors.WithMessage(err, "init HTTP Server")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stopChan)

	//HTTP Server go routine
	go func() {
		log.Printf("WhatFlix HTTP server is up and running on %s", config.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan <- err
			}
		}
	}()

	// WhatFlix HTTP Server Health Check
	go httphealthcheck.Check(healthCheckAddr, errChan)

	defer func() {
		log.Println("Gracefully Shutting Down WhatFlix HTTP server...")
		time.Sleep(5 * time.Second)
		httpServer.Shutdown(ctx)
		log.Println("WhatFlix HTTP Server is now Closed.")
		close(errChan)
		close(stopChan)
	}()

	select {
	case err := <-errChan:
		log.Printf("Fatal error WhatFlix Server: %v\n", err)
	case <-stopChan:
		log.Println("received shutdown/termination signal")
	case <-ctx.Done():
		cancel()
	}

	return nil
}

func handleError(err error) {
	log.Fatalf("%s", err)
}
