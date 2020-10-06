package jwtpkg

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Claims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWTToken generates a JWT token with the username and singed by the given secret key
func GenerateJWTToken(userName, jwtAccSecretKey string) (string, error) {
	claims := jwt.MapClaims{
		"username":  userName,
		"ExpiresAt": jwt.TimeFunc().Add(1 * time.Minute).Unix(),
		"IssuedAt":  jwt.TimeFunc().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtAccSecretKey))
}

// ValidateJWTToken validates a JWT token with the given key
func ValidateJWTToken(jwtKey, tokenString string) (*Claims, error) {
	clms := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, clms, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.WithMessage(err, "Invalid Signature")
		}
		return nil, errors.WithMessage(err, "Access Denied-Please check the access token")
	}

	if !tkn.Valid {
		return nil, errors.WithMessage(err, "Invalid Token")
	}
	return clms, nil
}
