package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
)

var mutex sync.Mutex
var entries logEntries

const logPath = "log.txt"

var tickCh = time.Tick(2 * time.Second)
var writeDelay = 2 * time.Second

func main() {
	err := run()
	if err != nil {
		handleError(err)
	}

}

func run() error {
	config, err := config.GetLogConfig()
	if err != nil {
		return errors.WithMessage(err, "get log config")
	}

	http.HandleFunc("/", whatflixLogEntry)

	f, err := os.Create(logPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create log file")
	}
	f.Close()

	srv := &http.Server{
		Addr: config.HTTPPort,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stopChan)

	//HTTP Log Server go routine
	go func() {
		log.Printf("WhatFlix Log server is up and running on port %s", config.HTTPPort)
		//if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan <- err
			}
		}
	}()

	go writeLog()

	defer func() {
		log.Println("Gracefully Shutting Down WhatFlix logging service...")
		time.Sleep(5 * time.Second)
		_ = srv.Shutdown(ctx)
		log.Println("WhatFlix logging service is now Closed.")
		close(errChan)
		close(stopChan)
	}()

	select {
	case err := <-errChan:
		log.Printf("Fatal error WhatFlix logging service: %v\n", err)
	case <-stopChan:
		log.Println("received shutdown/termination signal")
	case <-ctx.Done():
		cancel()
	}
	return err
}

func whatflixLogEntry(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var entry entity.LogEntry
	err := dec.Decode(&entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	entries = append(entries, entry)
	mutex.Unlock()
}

func writeLog() {
	for range tickCh {
		mutex.Lock()

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			log.Println(err)
			continue
		}
		targetTime := time.Now().Add(-writeDelay)
		sort.Sort(entries)
		for i, entry := range entries {
			if entry.Timestamp.Before(targetTime) {
				_, err := logFile.WriteString(writeEntry(entry))
				if err != nil {
					fmt.Println(err)
				}

				if i == len(entries)-1 {
					entries = logEntries{}
				}

			} else {
				entries = entries[i:]
				break
			}
		}

		logFile.Close()

		mutex.Unlock()
	}
}

func writeEntry(entry entity.LogEntry) string {

	return fmt.Sprintf("%v;%v;%v;%v\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Level, entry.Source, entry.Message)
}

func handleError(err error) {
	log.Fatalf("%s", err)
}
