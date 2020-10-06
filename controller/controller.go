package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/logservice/loghelper"
	"github.com/whatflix/pkg/httperrors"
)

func newHTTPHandler(cfg *config.Config) http.Handler {
	h := mux.NewRouter()
	//HeartBeat Check By Load Balancer
	h.NewRoute().
		Methods(http.MethodGet).
		//Path("/movies/ping").
		Path("/ping").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Pong"))
		})
	h.NewRoute().
		Methods(http.MethodPost).
		Path("/movies/signin").
		Handler(&authHTTPHandler{
			cfg:                     cfg,
			signinCredentialManager: &signinCredentialManagerImpl{},
		})
	h.NewRoute().
		Methods(http.MethodGet).
		Path("/movies/users").
		Handler(&recommendHTTPHandler{
			cfg:                    cfg,
			userPreferencesManager: &userPreferencesManagerImpl{},
			creditsManager:         &creditsManagerImpl{},
			moviesManager:          &moviesManagerImpl{},
		})
	h.NewRoute().
		Methods(http.MethodGet).
		Path("/movies/user/{$userID:[0-9]+}/search").
		Handler(&searchHTTPHandler{
			cfg:                    cfg,
			userPreferencesManager: &userPreferencesManagerImpl{},
			creditsManager:         &creditsManagerImpl{},
			moviesManager:          &moviesManagerImpl{},
		})
	return h
}

func Startup(cfg *config.Config) http.Handler {
	return newHTTPHandler(cfg)
}

func handleHTTP(w http.ResponseWriter, req *http.Request, src, logURL string, f func(w http.ResponseWriter, src, logURL string, req *http.Request) ([]byte, *httperrors.HTTPErr)) {
	response, httpERR := f(w, src, logURL, req)
	if httpERR != nil {
		http.Error(w, httpERR.Error, httpERR.StatusCode)
		go loghelper.WriteEntry(logURL, &entity.LogEntry{
			Level:     entity.LogLevelError,
			Timestamp: time.Now(),
			Source:    src,
			Message:   httpERR.Error,
		})
		err, _ := json.MarshalIndent(httpERR, "", "")
		log.Printf("Error: %s", string(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(response)
	if err != nil {
		log.Printf("error while writing http respose :%v", err)
	}
}
