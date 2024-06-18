package json

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	gohungry "github.com/guhungry/gohungry/http"
)

// MockHTTPClient is a mock implementation of the HTTPClient interface for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do executes the mocked HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// Mock response data structure
type MockResponse struct {
	Message string `json:"message"`
}

func TestGet(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(`{"message":"success"}`))
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		},
	}
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	response, err := Get[MockResponse]("http://example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Message != "success" {
		t.Errorf("Expected message to be 'success', got %s", response.Message)
	}
}

func TestPost(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(`{"message":"created"}`))
			return &http.Response{
				StatusCode: 201,
				Body:       body,
			}, nil
		},
	}
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	requestBody := map[string]string{"key": "value"}
	response, err := Post[MockResponse]("http://example.com", requestBody)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Message != "created" {
		t.Errorf("Expected message to be 'created', got %s", response.Message)
	}
}

func TestGetWithError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	_, err := Get[MockResponse]("http://example.com")
	if err == nil || err.Error() != "network error" {
		t.Fatalf("Expected network error, got %v", err)
	}
}
