package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/internal/config"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh    = make(chan *webRequest)
	registerCh   = make(chan string)
	unregisterCh = make(chan string)
	heartbeatCh  = time.Tick(50 * time.Second)
)

var (
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		//DisableCompression: false,
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
	config, err := config.GetLoadBalConfig()
	if err != nil {
		return errors.WithMessage(err, "get loadbalancer config")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		requestCh <- &webRequest{r: r, w: w, doneCh: doneCh}
		<-doneCh
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go processRequests(config.LogServiceURL)

	errChan1 := make(chan error)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stopChan)

	srv1 := &http.Server{
		Addr: config.HTTPPort,
	}

	//HTTP Load Balancer Server go routine
	go func() {
		log.Printf("WhatFlix Load Balancer Server is up and running on %s", config.HTTPPort)
		//if err := srv1.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv1.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan1 <- err
			}
		}
	}()

	errChan2 := make(chan error)

	srv2 := &http.Server{
		Addr:    config.IntHTTPPort,
		Handler: new(appserverHandler),
		/*Handler: &middleware.TimeoutMiddleware{
			Next: &middleware.GzipMiddleware{
				Next: new(appserverHandler),
			},
		},*/
	}

	//HTTP Load Balancer Registration Server go routine
	go func() {
		log.Printf("WhatFlix Load Balancer Registration Server is up and running on %s", config.IntHTTPPort)
		//if err := srv2.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv2.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan2 <- err
			}
		}
	}()

	defer func() {
		log.Println("Gracefully Shutting Down WhatFlix Load Balancer Registration and Load Balancer servers...")
		time.Sleep(5 * time.Second)
		_ = srv1.Shutdown(ctx)
		_ = srv2.Shutdown(ctx)
		log.Println("WhatFlix Load Balancer Registration and Load Balacer Servers are now Closed.")
		close(errChan1)
		close(errChan2)
		close(stopChan)
	}()

	select {
	case <-errChan1:
		log.Println("WhatFlix Load Balancer Server is down")
	case <-errChan2:
		log.Println("WhatFlix Load Balancer Registration Server is down")
	case <-stopChan:
		log.Println("received shutdown/termination signal")
	case <-ctx.Done():
		cancel()
	}
	return nil
}

type appserverHandler struct{}

func (h *appserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	port := r.URL.Query().Get("port")
	// log.Println("ip :", ip)
	// log.Println("port :", port)
	// log.Println("URL Path :", r.URL.Path)
	switch r.URL.Path {
	case "/register":
		registerCh <- ip + port
	case "/unregister":
		unregisterCh <- ip + port
	}
}

func handleError(err error) {
	log.Fatalf("%s", err)
}
