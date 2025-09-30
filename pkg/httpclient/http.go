package httpclient

import (
	"net/http"
	"strings"
)

type ResponseMetadata struct {
	Err        error
	StatusCode int
}

// MakeRequest make a get request if there no socket listening return a nil instance
func MakeRequest(pathUrl string) *ResponseMetadata {
	resp, err := http.Get(pathUrl)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil
		}
		return &ResponseMetadata{
			Err: err,
		}
	}
	defer resp.Body.Close()
	return &ResponseMetadata{
		Err:        nil,
		StatusCode: resp.StatusCode,
	}
}
