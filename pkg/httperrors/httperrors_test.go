package httperrors

import (
	"net/http"
	"testing"
)

const (
	source  = "source"
	message = "testMessage"
	err     = "testError"
)

func TestBadRequestError(t *testing.T) {
	wantSource := "source"
	wantMessage := "testMessage"
	wantErr := "testError"
	wantStatusCode := http.StatusBadRequest

	got := NewBadRequestError(source, message, err)
	if got.Source != wantSource {
		t.Fatalf("unexpected source got %s\n want %s", got.Source, wantSource)
	}
	if got.Message != wantMessage {
		t.Fatalf("unexpected message got %s\n want %s", got.Message, wantMessage)
	}
	if got.Error != wantErr {
		t.Fatalf("unexpected err got %s\n want %s", got.Error, wantErr)
	}
	if got.StatusCode != wantStatusCode {
		t.Fatalf("unexpected statuscode got %d\n want %d", got.StatusCode, wantStatusCode)
	}
}

func TestForbiddenError(t *testing.T) {
	wantStatusCode := http.StatusForbidden

	got := NewForbiddenError(source, message, err)
	if got.StatusCode != wantStatusCode {
		t.Fatalf("unexpected statuscode got %d\n want %d", got.StatusCode, wantStatusCode)
	}
}

func TestInternalServerError(t *testing.T) {
	wantStatusCode := http.StatusInternalServerError

	got := NewInternalServerError(source, message, err)
	if got.StatusCode != wantStatusCode {
		t.Fatalf("unexpected statuscode got %d\n want %d", got.StatusCode, wantStatusCode)
	}
}

func TestUnprocessableEntityError(t *testing.T) {
	wantStatusCode := http.StatusUnprocessableEntity

	got := NewUnprocessableEntityError(source, message, err)
	if got.StatusCode != wantStatusCode {
		t.Fatalf("unexpected statuscode got %d\n want %d", got.StatusCode, wantStatusCode)
	}
}
