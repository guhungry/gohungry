// Package httptest provides utilities for testing HTTP clients and servers.
// It includes dummy serializers and parsers, as well as mock implementations
// of the HTTPClient interface to simulate various HTTP responses.
package httptest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// DummyRequestBodySerializer is a dummy implementation for serializing request bodies during testing.
// It returns a static byte slice regardless of input.
func DummyRequestBodySerializer(body any) ([]byte, error) {
	return []byte("dummy body"), nil
}

// DummyResponseBodyParser is a dummy implementation for parsing response bodies during testing.
// It decodes the JSON content of the reader into a Response object.
func DummyResponseBodyParser[Response any](reader io.ReadCloser) (*Response, error) {
	var result Response
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MockHTTPClient is a mock implementation of the HTTPClient interface for testing.
// The behavior of the Do method is controlled by the DoFunc field.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do executes the mocked HTTP request by calling the DoFunc field.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockHTTPClientSuccess creates a MockHTTPClient that always returns a successful HTTP response with the given status code and response body.
func MockHTTPClientSuccess(statusCode int, response string) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(response))
			return &http.Response{
				StatusCode: statusCode,
				Body:       body,
			}, nil
		},
	}
}

// MockHTTPClientError creates a MockHTTPClient that always returns an error with the given message.
func MockHTTPClientError(message string) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New(message)
		},
	}
}
