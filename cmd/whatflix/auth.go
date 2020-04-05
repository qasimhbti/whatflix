package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// CreateTokenEndpoint validates the user credentials
func (a *App) createTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	var user User
	user.UserName = username
	err = user.get(a.DB)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	a.Role = user.Role
	if password == user.password {
		// Create a claims map
		claims := jwt.MapClaims{
			"username":  username,
			"ExpiresAt": 15000,
			"IssuedAt":  time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		response := Response{Token: tokenString, Status: "success"}
		respondWithJSON(w, http.StatusOK, response)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}
}

// ProtectedEndpoint returns the date and time
func (a *App) protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Access Denied; Please check the access token")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g.[]byte("my_secret_key")
		return secretKey, nil
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Access Denied2; Please check the access token")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// If token is valid
		response := make(map[string]string)
		response["time"] = time.Now().String()
		response["user"] = claims["username"].(string)
		respondWithJSON(w, http.StatusOK, response)
	} else {
		respondWithError(w, http.StatusForbidden, err.Error())
	}
}
