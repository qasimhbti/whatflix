package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/pkg/httperrors"
)

type searchHTTPHandler struct {
	cfg                    *config.Config
	userPreferencesManager interface {
		get(userID int32) (*entity.UserPreferences, error)
	}
	moviesManager interface {
		get(title string, prefLangSF []string) ([]*entity.MoviesCollRecord, error)
	}
	creditsManager interface {
		get(text *entity.SearchText) ([]*entity.CreditsData, error)
	}
}

func (h *searchHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handleHTTP(w, req, "Whatflix-Search", h.cfg.LogServiceURL, h.handle)
}

func (h *searchHTTPHandler) handle(w http.ResponseWriter, src, logURL string, req *http.Request) ([]byte, *httperrors.HTTPErr) {
	_, err := protectedEndpoint(h.cfg.JWTAccessSecretKey, req)
	if err != nil {
		err := httperrors.NewForbiddenError(src, "Invalid Login Credentials", err.Error())
		return nil, err
	}

	userID, reqTexts, err := h.getRequestData(req)
	if err != nil {
		err := httperrors.NewBadRequestError(src, "get request data", err.Error())
		return nil, err
	}

	userPrefs, err := h.userPreferencesManager.get(int32(userID))
	if err != nil {
		err := httperrors.NewUnprocessableEntityError(src, "get user preferences", err.Error())
		return nil, err
	}

	//Case : 1
	//    ** UserID does not present in User Preferences json
	if userPrefs == nil {

		response, err := h.processSearchTexts(reqTexts)
		if err != nil {
			err := httperrors.NewInternalServerError(src, "process search texts", err.Error())
			return nil, err
		}
		resp, _ := json.Marshal(map[string]interface{}{"Success": response})
		return resp, nil
	}

	//Case : 2
	//    ** UserID present in User Preferences json
	if userPrefs != nil {
		searchTexts := h.userPrefSearchTextGetter(reqTexts, userPrefs)

		//No Search text matches with users preferences actors or directors
		if len(searchTexts.Actors) == 0 && len(searchTexts.Directors) == 0 {
			response, err := h.processSearchTexts(reqTexts)
			if err != nil {
				err := httperrors.NewInternalServerError(src, "process search texts", err.Error())
				return nil, err
			}
			resp, _ := json.Marshal(map[string]interface{}{"Success": response})
			return resp, nil
		}

		//Search text matches with users preferences favorite actors or directors
		if len(searchTexts.Actors) != 0 || len(searchTexts.Directors) != 0 {
			creditsDatas, err := h.creditsManager.get(searchTexts)
			if err != nil {
				err := httperrors.NewInternalServerError(src, "credit collection", err.Error())
				return nil, err
			}

			if creditsDatas != nil {
				movies := removeDuplicate(creditsDatas)
				var userMoviesRecord []*entity.MoviesCollRecord
				for _, title := range movies {
					moviesRecord, err := h.moviesManager.get(title, userPrefs.PrefLangShortForm)
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
					resp, _ := json.Marshal(map[string]interface{}{"Success": movies})
					return resp, nil
				}

				response, err := h.processSearchTexts(reqTexts)
				if err != nil {
					err := httperrors.NewInternalServerError(src, "process search texts", err.Error())
					return nil, err
				}
				resp, _ := json.Marshal(map[string]interface{}{"Success": response})
				return resp, nil
			}
			resp, _ := json.Marshal(map[string]interface{}{"Success": "No movie found!!"})
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
	if queryText[0] == "" {
		return userID, []string{}, errors.New("invalid search text")
	}
	return userID, strings.Split(queryText[0], ","), nil
}

func (h *searchHTTPHandler) userPrefSearchTextGetter(texts []string, preferences *entity.UserPreferences) *entity.SearchText {
	var userSearchText entity.SearchText
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

func (h *searchHTTPHandler) processSearchTexts(reqTexts []string) ([]string, error) {
	var searchTexts entity.SearchText
	searchTexts.Texts = append(searchTexts.Texts, reqTexts...)

	creditsDatas, err := h.creditsManager.get(&searchTexts)
	if err != nil {
		return nil, errors.WithMessage(err, "credit collection")
	}
	if creditsDatas == nil {
		return []string{"No movie found!!"}, nil
	}

	var movies []string
	for _, data := range creditsDatas {
		movie := *data
		movies = append(movies, movie.Title)
	}

	sort.Strings(movies)
	return movies, nil
}
