package http

import "net/http"

// Client HTTP Client default to http.Client, replaced with mocks in unit test
var Client HTTPClient

func init() {
	Client = &http.Client{}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
