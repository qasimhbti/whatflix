package main

import (
	"log"
	"net/http"

	"github.com/whatflix/pkg/httperrors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func initHTTPServer(config *configs) (*http.Server, error) {
	DBClient, err := newMongoDBClientGetter(config.DBConString)
	if err != nil {
		return nil, errors.WithMessage(err, "mongo client")
	}
	//defer DBClient.Disconnect(ctx)

	return startHTTPServer(config.HTTPPort, DBClient), nil
}

func startHTTPServer(addr string, DBClient *mongo.Client) *http.Server {
	h := newHTTPHandler(DBClient)

	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

func newHTTPHandler(DBClient *mongo.Client) http.Handler {
	r := mux.NewRouter()
	r.NewRoute().
		Methods(http.MethodGet).
		Path("/movies/user/{$userID:[0-9]+}/search").
		Handler(&searchHTTPHandler{
			DBClient:               DBClient,
			userPreferencesManager: &userPreferencesManagerImpl{},
			creditsManager:         &creditsManagerImpl{},
			moviesManager:          &moviesManagerImpl{},
		})
	r.NewRoute().
		Methods(http.MethodGet).
		Path("/movies/users").
		Handler(&recommendHTTPHandler{
			DBClient:               DBClient,
			userPreferencesManager: &userPreferencesManagerImpl{},
			creditsManager:         &creditsManagerImpl{},
			moviesManager:          &moviesManagerImpl{},
		})
	r.NewRoute().
		Methods(http.MethodGet).
		Path("/movies/ping").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Pong"))
		})
	return r
}

func handleHTTP(w http.ResponseWriter, req *http.Request, name string, f func(req *http.Request) ([]byte, error)) {
	response, err := f(req)
	if err != nil {
		err = errors.WithMessage(err, name)
		err = errors.WithMessage(err, "handle http")
		code, text := httperrors.GetCodeText(err)
		http.Error(w, text, code)
		log.Printf("Error : %s\n%+v", err, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
