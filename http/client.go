package http

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client HTTP Client default to http.Client, replaced with mocks in unit test
var Client HTTPClient

func init() {
	newDefaultClient()
}

func newDefaultClient() {
	Client = &http.Client{}
}