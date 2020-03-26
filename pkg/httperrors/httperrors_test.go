package httperrors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

func TestGetCodeText(t *testing.T) {
	for _, tc := range []struct {
		name         string
		err          error
		expectedCode int
		expectedText string
	}{
		{
			name: "WithCode",
			err: WithCode(
				errors.New("test"),
				http.StatusBadRequest,
			),
			expectedCode: http.StatusBadRequest,
			expectedText: "test",
		},
		{
			name: "Wrapped",
			err: errors.WithMessage(
				WithCode(
					errors.New("test 2"),
					http.StatusBadRequest,
				),
				"test 1",
			),
			expectedCode: http.StatusBadRequest,
			expectedText: "test 2",
		},
		{
			name:         "Unknown",
			err:          errors.New("test"),
			expectedCode: http.StatusInternalServerError,
			expectedText: http.StatusText(http.StatusInternalServerError),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			code, text := GetCodeText(tc.err)
			if code != tc.expectedCode {
				t.Fatalf("unexpected code: got %d, want %d", code, tc.expectedCode)
			}
			if text != tc.expectedText {
				t.Fatalf("unexpected text: got %q, want %q", text, tc.expectedText)
			}
		})
	}
}

func TestWithCodeNil(t *testing.T) {
	err := WithCode(nil, http.StatusBadRequest)
	if err != nil {
		t.Fatal("not nil")
	}
}

func TestError(t *testing.T) {
	err := WithCode(errors.New("test"), http.StatusBadRequest)
	_ = err.Error()
}

func TestCause(t *testing.T) {
	cer := errors.New("1")
	err := WithCode(cer, http.StatusBadRequest)
	if errors.Cause(err) != cer {
		t.Fatalf("unexpected cause: got %#v, want %#v", errors.Cause(err), cer)
	}
	_ = fmt.Sprintf("%s", err)
	_ = fmt.Sprintf("%v", err)
	_ = fmt.Sprintf("%+v", err)
}
