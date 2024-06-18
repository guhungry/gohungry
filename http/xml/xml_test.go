package xml

import (
	"bytes"
	"encoding/xml"
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
	Message string `xml:"message"`
}

func TestGet(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(`<MockResponse><message>success</message></MockResponse>`))
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
			body := io.NopCloser(bytes.NewBufferString(`<MockResponse><message>created</message></MockResponse>`))
			return &http.Response{
				StatusCode: 201,
				Body:       body,
			}, nil
		},
	}
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	requestBody := struct {
		XMLName xml.Name `xml:"Key"`
		Value   string   `xml:"Value"`
	}{Value: "value"}
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

func TestPostWithInvalidXML(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(`invalid xml`))
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		},
	}
	gohungry.SetHTTPClient(mockClient)
	defer gohungry.ResetHTTPClient()

	requestBody := struct {
		Key string `xml:"key"`
	}{Key: "value"}
	_, err := Post[MockResponse]("http://example.com", requestBody)
	if err == nil {
		t.Fatalf("Expected XML decoding error, got no error")
	}
}
