package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/controller"
	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/internal/version"
	"github.com/whatflix/logservice/loghelper"
	"github.com/whatflix/middleware"
	"github.com/whatflix/model"
	"github.com/whatflix/pkg/mongodb"
	"github.com/whatflix/utils"
)

//const healthCheckAddr = "/movies/health-check"

var (
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: false,
	}
)

func init() {
	http.DefaultClient = &http.Client{Transport: &transport}
}

func main() {
	err := run()
	if err != nil {
		handleError(err)
	}
}

func run() error {
	utils.InitLog()
	config, err := config.GetConfig()
	if err != nil {
		return errors.WithMessage(err, "get config")
	}
	utils.LogStart(version.VersionString("Whatflix-web"), config.Environment)

	err = connectToDatabase(config.DBConString, config.DBName)
	if err != nil {
		return errors.WithMessage(err, "DB connection")
	}

	//HeartBeat Check By Load Balancer
	//http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//})

	handler := controller.Startup(config)

	srv := &http.Server{
		//Addr: config.HTTPHost + config.HTTPPort,
		Addr: config.HTTPPort,
		Handler: &middleware.TimeoutMiddleware{
			Next: &middleware.GzipMiddleware{
				Next: handler,
			},
		},
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
		log.Printf("WhatFlix HTTP server is up and running on port %s", config.HTTPPort)
		//if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan <- err
			}
		}
	}()

	// WhatFlix HTTP Server Health Check
	//go httphealthcheck.Check(healthCheckAddr, config.HTTPPort, errChan)

	//Service Discovery By Load Balancer
	go func() {
		time.Sleep(1 * time.Second)
		loghelper.WriteEntry(config.LogServiceURL, &entity.LogEntry{
			Level:     entity.LogLevelInfo,
			Timestamp: time.Now(),
			Source:    "Whatflix app server",
			Message:   "Registering with Load Balancer",
		})
		_, err = http.Get(config.LoadBalancerURL + "/register?port=" + config.HTTPPort)
		if err != nil {
			errChan <- err
		}
	}()

	defer func() {
		log.Println("Gracefully Shutting Down WhatFlix HTTP server...")
		time.Sleep(5 * time.Second)
		_ = srv.Shutdown(ctx)
		log.Println("WhatFlix HTTP Server is now Closed.")

		//Service Termination By Load Balancer
		go loghelper.WriteEntry(config.LogServiceURL, &entity.LogEntry{
			Level:     entity.LogLevelInfo,
			Timestamp: time.Now(),
			Source:    "Whatflix app server",
			Message:   "Unregistering with Load Balancer",
		})
		_, _ = http.Get(config.LoadBalancerURL + "/unregister?port=" + config.HTTPPort)
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

func connectToDatabase(cs, dbName string) error {
	DBClient, err := mongodb.NewMongoDBClientGetter(cs)
	if err != nil {
		return errors.WithMessage(err, "mongo client")
	}

	db := mongodb.GetmgoDB(DBClient, dbName)
	model.SetDatabase(db)
	return nil
}

func handleError(err error) {
	log.Fatalf("%s", err)
}
