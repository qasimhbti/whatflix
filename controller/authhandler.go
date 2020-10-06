package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/whatflix/entity"
	"github.com/whatflix/internal/config"
	"github.com/whatflix/internal/jwtpkg"
	"github.com/whatflix/pkg/httperrors"
)

type authHTTPHandler struct {
	cfg                     *config.Config
	signinCredentialManager interface {
		get(userName string) (*entity.SigninCred, error)
	}
}

func (h *authHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handleHTTP(w, req, "Whatflix-Signin", h.cfg.LogServiceURL, h.handle)
}

type signinResponse struct {
	Key         string `json:"key"`
	TokenString string `json:"tokenstring"`
	Status      string `json:"status"`
	//ExpiresAt   int64  `json:"expiresat"`
}

func (h *authHTTPHandler) handle(w http.ResponseWriter, src, logURL string, r *http.Request) ([]byte, *httperrors.HTTPErr) {
	userCreds, err := h.getUserCredentials(r)
	if err != nil {
		return nil, httperrors.NewBadRequestError(src, "get user creds", err.Error())
	}

	//expiration time of token
	//expirationTime := time.Now().Add(1 * time.Minute)
	tokenString, err := jwtpkg.GenerateJWTToken(userCreds.UserName, h.cfg.JWTAccessSecretKey)
	if err != nil {
		return nil, httperrors.NewInternalServerError(src, "generating signed token", err.Error())
	}

	resp, err := json.Marshal(&signinResponse{
		Key:         "Access-Token",
		TokenString: tokenString,
		Status:      "Successfully Generated",
		//ExpiresAt:   expirationTime.Unix(),
	})
	if err != nil {
		return nil, httperrors.NewUnprocessableEntityError(src, "JSON Marshal", err.Error())
	}

	http.SetCookie(w, &http.Cookie{
		Name:  src,
		Value: tokenString,
		//Expires: expirationTime,
	})

	return resp, nil
}

func (h *authHTTPHandler) getUserCredentials(r *http.Request) (*entity.SigninCred, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	userName := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	if userName == "" || password == "" {
		return nil, errors.New("username or password is empty")
	}

	signinCreds, err := h.signinCredentialManager.get(userName)
	if err != nil {
		return nil, err
	}

	if password != signinCreds.Password {
		return nil, errors.New("invalid credentials")
	}

	return signinCreds, nil
}

// ProtectedEndpoint returns the date and time
func protectedEndpoint(jwtKey string, r *http.Request) (*jwtpkg.Claims, error) {
	cookie, err := r.Cookie("Whatflix-Signin")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			log.Println("Cookie not set")
		}
		// For any other type of error, return a bad request status
		return nil, errors.WithMessage(err, "Invalid Request")
	}

	return jwtpkg.ValidateJWTToken(jwtKey, cookie.Value)
}
