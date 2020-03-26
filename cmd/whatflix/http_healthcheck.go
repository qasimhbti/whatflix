package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

var isHealth bool = true

func whatflixHealthCheck(errChan chan error) {
	for {
		if isHealth {
			time.Sleep(time.Second * 2)
			// Create a Health Check Request
			req, err := http.NewRequest(http.MethodGet, "/movies/health-check", nil)
			if err != nil {
				log.Println("error while creating health check request :", err)
				continue
			}

			// ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(healthCheckHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code == http.StatusOK {
				log.Println("Health Check : WhatFlix HTTP Server is up and running")
			} else {
				isHealth = false
				log.Println("Health Check : WhatFlix HTTP server is NOT SERVING")
				errChan <- errors.New("Health Check Faild")
			}
			time.Sleep(time.Second * 58)
		}
	}
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
