package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock implementation of the HTTPClient interface for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do executes the mocked HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestDoRequest tests the DoRequest function.
func TestDoRequest(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewBufferString(`{"message":"success"}`))
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		},
	}
	SetHTTPClient(mockClient)
	defer ResetHTTPClient()

	tests := []struct {
		name           string
		requestInfo    *RequestInfo[map[string]interface{}]
		expectedBody   string
		expectedHeader string
		expectedError  error
	}{
		{
			name: "Successful Request",
			requestInfo: &RequestInfo[map[string]interface{}]{
				method:         MethodPost,
				url:            "https://example.com",
				body:           map[string]string{"key": "value"},
				bodySerializer: dummyRequestBodySerializer,
				responseParser: dummyResponseBodyParser[map[string]interface{}],
				headers:        Headers{HeaderContentType: "application/json"},
			},
			expectedBody:   "dummy body",
			expectedHeader: "application/json",
			expectedError:  nil,
		},
		{
			name: "Error in Serialization",
			requestInfo: &RequestInfo[map[string]interface{}]{
				method:         MethodPost,
				url:            "https://example.com",
				body:           map[string]string{"key": "value"},
				bodySerializer: func(body any) ([]byte, error) { return nil, io.ErrUnexpectedEOF },
				responseParser: dummyResponseBodyParser[map[string]interface{}],
				headers:        Headers{HeaderContentType: "application/json"},
			},
			expectedBody:   "",
			expectedHeader: "",
			expectedError:  io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := DoRequest(tt.requestInfo)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("DoRequest() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError != nil {
				return
			}

			// Check request body
			if (*response)["message"] != "success" {
				t.Errorf("Expected response message to be 'success', got %v, %+v", (*response)["message"], nil)
			}

			// Check headers
			req, _ := http.NewRequest(tt.requestInfo.method, tt.requestInfo.url, nil)
			setHeaders(req, tt.requestInfo.headers)
			if req.Header.Get(HeaderContentType) != tt.expectedHeader {
				t.Errorf("Expected header %s to be %s, got %s", HeaderContentType, tt.expectedHeader, req.Header.Get(HeaderContentType))
			}
		})
	}
}

// TestSetAuth tests the setAuth function for both basic and bearer authentication.
func TestSetAuth(t *testing.T) {
	tests := []struct {
		name               string
		authType           AuthType
		authCredentials    AuthCredentials
		expectedAuthHeader string
	}{
		{
			name:     "Basic Authentication",
			authType: AuthTypeBasic,
			authCredentials: AuthCredentials{
				username: "user",
				password: "pass",
			},
			expectedAuthHeader: string(AuthTypeBasic) + " dXNlcjpwYXNz",
		},
		{
			name:     "Bearer Authentication",
			authType: AuthTypeBearer,
			authCredentials: AuthCredentials{
				token: "some-token",
			},
			expectedAuthHeader: string(AuthTypeBearer) + " some-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			requestInfo := &RequestInfo[any]{
				authType:        tt.authType,
				authCredentials: tt.authCredentials,
			}
			setAuth(req, requestInfo)

			authHeader := req.Header.Get(HeaderAuthorization)
			if authHeader != tt.expectedAuthHeader {
				t.Errorf("Expected Authorization header '%s', got '%s'", tt.expectedAuthHeader, authHeader)
			}
		})
	}
}

// TestToBodyReader tests the toBodyReader function.
func TestToBodyReader(t *testing.T) {
	tests := []struct {
		name          string
		body          any
		serializer    RequestBodySerializer
		expectedBody  string
		expectedError error
	}{
		{
			name:          "Nil Body",
			body:          nil,
			serializer:    dummyRequestBodySerializer,
			expectedBody:  "",
			expectedError: nil,
		},
		{
			name:          "Valid Body",
			body:          map[string]string{"key": "value"},
			serializer:    dummyRequestBodySerializer,
			expectedBody:  "dummy body",
			expectedError: nil,
		},
		{
			name:          "Serialization Error",
			body:          map[string]string{"key": "value"},
			serializer:    func(body any) ([]byte, error) { return nil, io.ErrUnexpectedEOF },
			expectedBody:  "",
			expectedError: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := toBodyReader(tt.body, tt.serializer)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("toBodyReader() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError != nil {
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(reader)
			if buf.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, buf.String())
			}
		})
	}
}
