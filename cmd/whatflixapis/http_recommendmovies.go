package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/pkg/errors"
	"github.com/whatflix/pkg/httperrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type recommendHTTPHandler struct {
	DBClient               *mongo.Client
	userPreferencesManager interface {
		getAll(db *mongo.Database) ([]*userPreferences, error)
	}
	creditsManager interface {
		get(text *searchText, db *mongo.Database) ([]*creditsData, error)
	}
	moviesManager interface {
		get(title string, prefLangSF []string, db *mongo.Database) ([]*moviesCollRecord, error)
	}
}

type userRecommendation struct {
	UserID            int      `json:"user"`
	RecommendedMovies []string `json:"movies"`
}

func (h *recommendHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handleHTTP(w, req, "Whatflix-Recommendation", h.handle)
}

func (h *recommendHTTPHandler) handle(req *http.Request) ([]byte, error) {
	db := getmgoDB(h.DBClient)
	userPrefs, err := h.userPreferencesManager.getAll(db)
	if err != nil {
		return nil, errors.WithMessage(
			httperrors.WithCode(err, http.StatusInternalServerError),
			"user preferences",
		)
	}

	if userPrefs == nil {
		resp := responseWithJSON("Success", []string{"No user found!!"})
		return resp, nil
	}

	var uRecommd []*userRecommendation
	for _, user := range userPrefs {
		var u userRecommendation
		var userMoviesRecord []*moviesCollRecord
		var text searchText
		for _, actor := range user.FavouriteActors {
			text.Actors = append(text.Actors, actor)
		}

		for _, director := range user.FavouriteDirectors {
			text.Directors = append(text.Directors, director)
		}

		creditsDatas, err := h.creditsManager.get(&text, db)
		if err != nil {
			log.Printf("error searching text : %+v in credits collection", text)
			continue
		}

		titles := removeDuplicate(creditsDatas)
		for _, title := range titles {
			moviesRecord, err := h.moviesManager.get(title, user.PrefLangShortForm, db)
			if err != nil {
				log.Printf("error retriving title :%s from movies collection", title)
				continue
			}
			userMoviesRecord = append(userMoviesRecord, moviesRecord...)
		}

		sort.Sort(ByVoteAverage(userMoviesRecord))
		count := 0
		for _, value := range userMoviesRecord {
			if count == 3 {
				break
			}
			u.RecommendedMovies = append(u.RecommendedMovies, value.Title)
			count++
		}
		u.UserID = user.UserID
		sort.Strings(u.RecommendedMovies)
		uRecommd = append(uRecommd, &u)
	}

	resp, err := json.Marshal(uRecommd)
	if err != nil {
		return nil, errors.WithMessage(
			httperrors.WithCode(err, http.StatusInternalServerError),
			"JSON Marshal",
		)
	}
	return resp, nil
}
