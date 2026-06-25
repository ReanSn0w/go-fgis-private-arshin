package client

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
}

func NewHTTPError(resp *http.Response) *HTTPError {
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("arshin: unexpected HTTP status %s", e.Status)
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("arshin: API status %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("arshin: API status %d", e.Status)
}

type UnexpectedContentTypeError struct {
	StatusCode  int
	ContentType string
}

func (e *UnexpectedContentTypeError) Error() string {
	return fmt.Sprintf("arshin: expected JSON response, got %q with status %d", e.ContentType, e.StatusCode)
}
