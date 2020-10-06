package jwtpkg

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	mockToken        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MDk1MDQ1NTYsIklzc3VlZEF0IjoxNjA5NTA0NDk2LCJ1c2VybmFtZSI6ImFkbWluIn0.a7JUyyRGRQz7_dxPYNaYkkBTU8C7GiEvhSGLEgald84"
	mockJWTSecretKey = "jdnfksdmfksd"
)

func init() {
	// Mock JWT time
	jwt.TimeFunc = func() time.Time {
		return time.Date(2021, 01, 01, 12, 34, 56, 0, time.UTC)
	}
}

func TestGenerateJWTToken(t *testing.T) {
	userName := "admin"

	token, err := GenerateJWTToken(userName, mockJWTSecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if token != mockToken {
		t.Fatalf("unexpected token: got %s \nwant %s", token, mockToken)
	}
}

func TestValidateJWTToken(t *testing.T) {

	_, err := ValidateJWTToken(mockJWTSecretKey, mockToken)
	if err != nil {
		t.Fatal(err)
	}
}
