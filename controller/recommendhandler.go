package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/pkg/httperrors"
)

type recommendHTTPHandler struct {
	cfg                    *config.Config
	userPreferencesManager interface {
		getAll() ([]*entity.UserPreferences, error)
	}
	creditsManager interface {
		get(text *entity.SearchText) ([]*entity.CreditsData, error)
	}
	moviesManager interface {
		get(title string, prefLangSF []string) ([]*entity.MoviesCollRecord, error)
	}
}

type moviesRecommendation struct {
	UserID            int32    `json:"user"`
	RecommendedMovies []string `json:"movies"`
}

func (h *recommendHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handleHTTP(w, req, "Whatflix-Recommendation", h.cfg.LogServiceURL, h.handle)
}

func (h *recommendHTTPHandler) handle(w http.ResponseWriter, src, logURL string, req *http.Request) ([]byte, *httperrors.HTTPErr) {
	_, err := protectedEndpoint(h.cfg.JWTAccessSecretKey, req)
	if err != nil {
		err := httperrors.NewForbiddenError(src, "Invalid Login Credentials", err.Error())
		return nil, err
	}

	cacheKey := req.URL.RequestURI()
	resp, ok := getFromCache(h.cfg.CacheServiceURL, cacheKey)
	if ok {
		_, _ = io.Copy(w, resp)
		resp.Close()
		return nil, nil
	}

	userPrefs, err := h.userPreferencesManager.getAll()
	if err != nil {
		err := httperrors.NewUnprocessableEntityError(src, "user preferences", err.Error())
		return nil, err
	}

	if userPrefs == nil {
		return []byte("Success: No user found!!"), nil
	}

	var moviesRecommd []*moviesRecommendation
	for _, user := range userPrefs {
		var mov moviesRecommendation
		var userMoviesRecord []*entity.MoviesCollRecord
		var text entity.SearchText

		text.Actors = append(text.Actors, user.FavouriteActors...)
		text.Directors = append(text.Directors, user.FavouriteDirectors...)

		creditsDatas, err := h.creditsManager.get(&text)
		if err != nil {
			log.Printf("error searching text : %+v in credits collection", text)
			continue
		}

		titles := removeDuplicate(creditsDatas)
		for _, title := range titles {
			moviesRecord, err := h.moviesManager.get(title, user.PrefLangShortForm)
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
			mov.RecommendedMovies = append(mov.RecommendedMovies, value.Title)
			count++
		}
		mov.UserID = user.UserID
		sort.Strings(mov.RecommendedMovies)
		moviesRecommd = append(moviesRecommd, &mov)
	}

	response, _ := json.Marshal(map[string]interface{}{"Succes": moviesRecommd})
	//go saveToCache(h.cfg.CacheServiceURL, cacheKey, int64(24*time.Hour), response)
	go saveToCache(h.cfg.CacheServiceURL, cacheKey, int64(70*time.Second), response)
	return response, nil
}

func removeDuplicate(items []*entity.CreditsData) []string {
	var key = make(map[string]bool)
	lists := []string{}
	for _, item := range items {
		entry := *item
		if _, value := key[entry.Title]; !value {
			key[entry.Title] = true
			lists = append(lists, entry.Title)
		}
	}
	return lists
}
