package httperrors

import "net/http"

type HTTPErr struct {
	Source     string `json:"source"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func NewBadRequestError(source, message, err string) *HTTPErr {
	return &HTTPErr{
		Source:     source,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Error:      err,
	}
}

func NewForbiddenError(source, message, err string) *HTTPErr {
	return &HTTPErr{
		Source:     source,
		Message:    message,
		StatusCode: http.StatusForbidden,
		Error:      err,
	}
}

func NewInternalServerError(source, message, err string) *HTTPErr {
	return &HTTPErr{
		Source:     source,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Error:      err,
	}
}

func NewUnprocessableEntityError(source, message, err string) *HTTPErr {
	return &HTTPErr{
		Source:     source,
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
		Error:      err,
	}
}
