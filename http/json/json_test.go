package json

import (
	"github.com/guhungry/gohungry/http/httptest"
	"testing"

	gohungry "github.com/guhungry/gohungry/http"
)

// Mock response data structure
type MockResponse struct {
	Message string `json:"message"`
}

func TestGet(t *testing.T) {
	mockClient := httptest.MockHTTPClientSuccess(200, `{"message":"success"}`)
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
	mockClient := httptest.MockHTTPClientSuccess(201, `{"message":"created"}`)
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
	mockClient := httptest.MockHTTPClientError("network error")
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	_, err := Get[MockResponse]("http://example.com")
	if err == nil || err.Error() != "network error" {
		t.Fatalf("Expected network error, got %v", err)
	}
}
