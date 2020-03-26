package httperrors

import (
	"fmt"
	"io"
	"net/http"
)

type base struct {
	cause error
	code  int
}

// WithCode annotates an error with an HTTP status code.
func WithCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &base{
		code:  code,
		cause: err,
	}
}

func (err *base) message() string {
	return fmt.Sprintf("HTTP code %d", err.code)
}

func (err *base) Error() string {
	return fmt.Sprintf("%s: %s", err.message(), err.cause)
}

func (err *base) Cause() error {
	return err.cause
}

func (err *base) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", err.cause)
			_, _ = io.WriteString(s, err.message())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, err.Error())
	}
}

// GetCodeText returns the HTTP code and text associated to an error.
func GetCodeText(err error) (code int, text string) {
	if err := get(err); err != nil {
		return err.code, err.cause.Error()
	}
	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
}

func get(err error) *base {
	for err != nil {
		if err, ok := err.(*base); ok {
			return err
		}
		err = getCause(err)
	}
	return nil
}

func getCause(err error) error {
	cer, ok := err.(interface {
		Cause() error
	})
	if !ok {
		return nil
	}
	return cer.Cause()
}
