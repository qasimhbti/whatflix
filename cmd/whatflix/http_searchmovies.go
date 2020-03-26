package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/whatflix/pkg/httperrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type searchHTTPHandler struct {
	DBClient               *mongo.Client
	userPreferencesManager interface {
		get(userID int, db *mongo.Database) (*userPreferences, error)
	}
	moviesManager interface {
		get(title string, prefLangSF []string, db *mongo.Database) ([]*moviesCollRecord, error)
	}
	creditsManager interface {
		get(text *searchText, db *mongo.Database) ([]*creditsData, error)
	}
}

func (h *searchHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handleHTTP(w, req, "Whatflix-Search", h.handle)
}

func (h *searchHTTPHandler) handle(req *http.Request) ([]byte, error) {
	userID, reqTexts, err := h.getRequestData(req)
	if err != nil {
		return nil, errors.WithMessage(
			httperrors.WithCode(err, http.StatusBadRequest),
			"get request data",
		)
	}

	db := getmgoDB(h.DBClient)
	userPrefs, err := h.userPreferencesManager.get(userID, db)
	if err != nil {
		return nil, errors.WithMessage(err, "get user preferences")
	}

	//Case : 1
	//    ** UserID does not present in User Preferences json AND
	//    ** No search Text in query
	if userPrefs == nil && len(reqTexts) == 0 {
		resp := responseWithJSON("Fail", []string{"UserID does not found and search text is blank"})
		return resp, nil
	}

	var searchTexts searchText
	//Case : 2
	//    ** UserID does not present in User Preferences json AND
	//    ** Search Text present in query
	if userPrefs == nil && len(reqTexts) != 0 {
		for _, reqText := range reqTexts {
			searchTexts.Texts = append(searchTexts.Texts, reqText)
		}

		creditsDatas, err := h.creditsManager.get(&searchTexts, db)
		if err != nil {
			return nil, errors.WithMessage(
				httperrors.WithCode(err, http.StatusInternalServerError),
				"credit collection",
			)
		}
		if creditsDatas == nil {
			resp := responseWithJSON("Success", []string{"No movie found!!"})
			return resp, nil
		}

		var movies []string
		for _, data := range creditsDatas {
			movie := *data
			movies = append(movies, movie.Title)
		}

		sort.Strings(movies)
		resp := responseWithJSON("Success", movies)
		return resp, nil
	}

	//Case : 3
	//    ** UserID present in User Preferences json AND
	//    ** No Search Text in query
	if userPrefs != nil && len(reqTexts) == 0 {
		for _, actor := range userPrefs.FavouriteActors {
			searchTexts.Actors = append(searchTexts.Actors, actor)
		}

		for _, director := range userPrefs.FavouriteDirectors {
			searchTexts.Directors = append(searchTexts.Directors, director)
		}

		creditsDatas, err := h.creditsManager.get(&searchTexts, db)
		if err != nil {
			return nil, errors.WithMessage(
				httperrors.WithCode(err, http.StatusInternalServerError),
				"credit collection",
			)
		}

		if creditsDatas == nil {
			resp := responseWithJSON("Success", []string{"No movie found!!"})
			return resp, nil
		}

		var userMoviesRecord []*moviesCollRecord
		titles := removeDuplicate(creditsDatas)
		for _, title := range titles {
			moviesRecord, err := h.moviesManager.get(title, userPrefs.PrefLangShortForm, db)
			if err != nil {
				log.Printf("error retriving title :%s from movies collection", title)
				continue
			}
			userMoviesRecord = append(userMoviesRecord, moviesRecord...)
		}

		var movies []string
		for _, data := range userMoviesRecord {
			movie := *data
			movies = append(movies, movie.Title)
		}

		sort.Strings(movies)
		resp := responseWithJSON("Success", movies)
		return resp, nil
	}

	//Case : 4
	//    ** UserID present in User Preferences json AND
	//    ** Search Text present in query
	if userPrefs != nil && len(reqTexts) != 0 {
		searchTexts := h.userPrefSearchTextGetter(reqTexts, userPrefs)
		if len(searchTexts.Actors) == 0 && len(searchTexts.Directors) == 0 {
			for _, reqText := range reqTexts {
				searchTexts.Texts = append(searchTexts.Texts, reqText)
			}

			creditsDatas, err := h.creditsManager.get(searchTexts, db)
			if err != nil {
				return nil, errors.WithMessage(
					httperrors.WithCode(err, http.StatusInternalServerError),
					"credit collection",
				)
			}
			if creditsDatas == nil {
				resp := responseWithJSON("Success", []string{"No movie found!!"})
				return resp, nil
			}

			var movies []string
			for _, data := range creditsDatas {
				movie := *data
				movies = append(movies, movie.Title)
			}

			sort.Strings(movies)
			resp := responseWithJSON("Success", movies)
			return resp, nil
		}

		if len(searchTexts.Actors) != 0 || len(searchTexts.Directors) != 0 {
			creditsDatas, err := h.creditsManager.get(searchTexts, db)
			if err != nil {
				return nil, errors.WithMessage(
					httperrors.WithCode(err, http.StatusInternalServerError),
					"credit collection",
				)
			}

			if creditsDatas != nil {
				movies := removeDuplicate(creditsDatas)
				var userMoviesRecord []*moviesCollRecord
				for _, title := range movies {
					moviesRecord, err := h.moviesManager.get(title, userPrefs.PrefLangShortForm, db)
					if err != nil {
						log.Printf("error retriving title :%s from movies collection", title)
						continue
					}
					userMoviesRecord = append(userMoviesRecord, moviesRecord...)
				}

				if userMoviesRecord != nil {
					var movies []string
					for _, data := range userMoviesRecord {
						movie := *data
						movies = append(movies, movie.Title)
					}

					sort.Strings(movies)
					resp := responseWithJSON("Success", movies)
					return resp, nil
				}
				//} else {
				for _, reqText := range reqTexts {
					searchTexts.Texts = append(searchTexts.Texts, reqText)
				}

				creditsDatas, err := h.creditsManager.get(searchTexts, db)
				if err != nil {
					return nil, errors.WithMessage(
						httperrors.WithCode(err, http.StatusInternalServerError),
						"credit collection",
					)
				}
				if creditsDatas == nil {
					resp := responseWithJSON("Success", []string{"No movie found!!"})
					return resp, nil
				}

				var movs []string
				for _, data := range creditsDatas {
					movie := *data
					movs = append(movs, movie.Title)
				}

				sort.Strings(movs)
				resp := responseWithJSON("Success", movs)
				return resp, nil
			}
			resp := responseWithJSON("Success", []string{"No movie found!!"})
			return resp, nil
		}
	}
	return nil, nil
}

func (h *searchHTTPHandler) getRequestData(req *http.Request) (int, []string, error) {
	vars := mux.Vars(req)
	userID, err := strconv.Atoi(vars["$userID"])
	if err != nil {
		return 0, nil, errors.WithMessage(err, "parse user id")
	}

	queryText := req.URL.Query()["text"]
	if queryText[0] != "" {
		return userID, strings.Split(queryText[0], ","), nil
	}
	return userID, nil, nil
}

type searchText struct {
	Texts     []string
	Actors    []string
	Directors []string
	Languages []string
}

func (h *searchHTTPHandler) userPrefSearchTextGetter(texts []string, preferences *userPreferences) *searchText {
	var userSearchText searchText
	for _, text := range texts {
		for _, prefActor := range preferences.FavouriteActors {
			if text == prefActor {
				userSearchText.Actors = append(userSearchText.Actors, text)
				//break
			}
		}

		for _, prefDirector := range preferences.FavouriteDirectors {
			if text == prefDirector {
				userSearchText.Directors = append(userSearchText.Directors, text)
				//break
			}
		}
	}
	return &userSearchText
}

type respMsg struct {
	Status  string   `json:"Status"`
	Message []string `json:"Movies"`
}

func responseWithJSON(status string, msg []string) []byte {
	response, _ := json.Marshal(&respMsg{
		Status:  status,
		Message: msg,
	})
	return response
}
