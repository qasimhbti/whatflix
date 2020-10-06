package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/internal/version"
	"github.com/whatflix/utils"
)

type cacheEntry struct {
	data       []byte
	expiration time.Time
}

var (
	cache  = make(map[string]*cacheEntry)
	mutex  = sync.RWMutex{}
	tickCh = time.Tick(50 * time.Second)
)

var maxAgeRexexp = regexp.MustCompile(`maxage=(\d+)`)

func main() {
	err := run()
	if err != nil {
		handleError(err)
	}
}

func run() error {
	utils.InitLog()
	config, err := config.GetCacheConfig()
	if err != nil {
		return errors.WithMessage(err, "get cache config")
	}
	utils.LogStart(version.VersionString("whatflix-cacheservice"), config.HTTPPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getFromCache(w, r)
		} else if r.Method == http.MethodPost {
			saveToCache(w, r)
		}
	})

	http.HandleFunc("/invalidate", invalidateEntry)

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

	//HTTP Cache Server go routine
	go func() {
		log.Printf("WhatFlix Cache server is up and running on %s", config.HTTPPort)
		//if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				errChan <- err
			}
		}
	}()

	go purgeCache()

	defer func() {
		log.Println("Gracefully Shutting Down WhatFlix caching service...")
		time.Sleep(5 * time.Second)
		_ = srv.Shutdown(ctx)
		log.Println("WhatFlix caching service is now Closed.")
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

func getFromCache(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	//key := r.URL.Query().Get("key")
	key := r.URL.RequestURI()

	log.Printf("Searching cache for %s", key)
	if entry, ok := cache[key]; ok {
		log.Println("found")
		_, _ = w.Write(entry.data)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	log.Println("not found")
}

func saveToCache(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	key := r.URL.RequestURI()
	cacheHeader := r.Header.Get("cache-control")

	log.Printf("Saving cache entry with key '%s' for %s seconds\n", key, cacheHeader)

	matches := maxAgeRexexp.FindStringSubmatch(cacheHeader)
	if len(matches) == 2 {
		dur, _ := strconv.Atoi(matches[1])
		data, _ := ioutil.ReadAll(r.Body)
		cache[key] = &cacheEntry{data: data, expiration: time.Now().Add(time.Duration(dur) * time.Second)}
	}
}

func invalidateEntry(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	//key := r.URL.Query().Get("key")
	key := r.URL.RequestURI()
	log.Printf("purging entry with key '%s'\n", key)
	delete(cache, key)
}

func purgeCache() {
	for range tickCh {
		mutex.Lock()

		for k := range cache {
			log.Printf("purging entry with key '%s'\n", k)
			delete(cache, k)
		}
		mutex.Unlock()
	}
}

func handleError(err error) {
	log.Fatalf("%s", err)
}
