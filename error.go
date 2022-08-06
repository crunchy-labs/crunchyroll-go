package crunchyroll

import (
	"fmt"
	"net/http"
)

// RequestError is an error interface which gets used whenever a crunchyroll delivers an error response.
type RequestError struct {
	error

	Response *http.Response
	Message  string
}

func (re *RequestError) Error() string {
	return fmt.Sprintf("error for endpoint %s (%d): %s", re.Response.Request.URL.String(), re.Response.StatusCode, re.Message)
}
