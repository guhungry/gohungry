package http

import "net/http"

// HTTPClient is an interface that wraps the Do method for making HTTP requests.
// This allows for easy replacement with mocks during unit testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the default HTTP client, which can be replaced with a mock in unit tests.
var Client HTTPClient = &http.Client{}

// SetHTTPClient replaces the default HTTP client with the provided one.
// This is useful for testing with a mock client.
func SetHTTPClient(client HTTPClient) {
	Client = client
}

// ResetHTTPClient restores the default HTTP client.
func ResetHTTPClient() {
	Client = &http.Client{}
}

func init() {
	ResetHTTPClient()
}
