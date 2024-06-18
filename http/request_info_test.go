package http

import (
	"io"
	"testing"
)

// Dummy implementations for testing
func dummyRequestBodySerializer(body any) ([]byte, error) {
	return []byte("dummy"), nil
}

func dummyResponseBodyParser[Response any](reader io.ReadCloser) (*Response, error) {
	var resp Response
	return &resp, nil
}

// TestNewRequestInfo tests the NewRequestInfo function with various options
func TestNewRequestInfo(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		body     any
		options  []RequestInfoOption[any]
		expected RequestInfo[any]
	}{
		{
			name:   "No Options",
			method: MethodGet,
			url:    "http://example.com",
			body:   nil,
			expected: RequestInfo[any]{
				method:         "GET",
				url:            "http://example.com",
				body:           nil,
				bodySerializer: dummyRequestBodySerializer,
				responseParser: dummyResponseBodyParser[any],
				headers:        make(Headers),
			},
		},
		{
			name:   "With Basic Auth",
			method: MethodPost,
			url:    "http://example.com",
			body:   "some body",
			options: []RequestInfoOption[any]{
				WithAuthBasic[any]("user", "pass"),
			},
			expected: RequestInfo[any]{
				method:         "POST",
				url:            "http://example.com",
				body:           "some body",
				bodySerializer: dummyRequestBodySerializer,
				responseParser: dummyResponseBodyParser[any],
				authType:       AuthTypeBasic,
				authCredentials: AuthCredentials{
					username: "user",
					password: "pass",
				},
				headers: make(Headers),
			},
		},
		{
			name:   "With Bearer Token",
			method: "PUT",
			url:    "http://example.com",
			body:   "some body",
			options: []RequestInfoOption[any]{
				WithAuthBearer[any]("some-token"),
			},
			expected: RequestInfo[any]{
				method:         "PUT",
				url:            "http://example.com",
				body:           "some body",
				bodySerializer: dummyRequestBodySerializer,
				responseParser: dummyResponseBodyParser[any],
				authType:       AuthTypeBearer,
				authCredentials: AuthCredentials{
					token: "some-token",
				},
				headers: make(Headers),
			},
		},
		{
			name:   "With Headers",
			method: "DELETE",
			url:    "http://example.com",
			body:   nil,
			options: []RequestInfoOption[any]{
				WithHeader[any]("Custom-Header", "value"),
				WithContentType[any]("application/json"),
				WithAccept[any]("application/json"),
			},
			expected: RequestInfo[any]{
				method:         "DELETE",
				url:            "http://example.com",
				body:           nil,
				bodySerializer: dummyRequestBodySerializer,
				responseParser: dummyResponseBodyParser[any],
				headers: Headers{
					"Custom-Header":   "value",
					HeaderContentType: "application/json",
					HeaderAccept:      "application/json",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewRequestInfo(tt.method, tt.url, tt.body, dummyRequestBodySerializer, dummyResponseBodyParser[any], tt.options...)
			if !equalRequestInfo(result, &tt.expected) {
				t.Errorf("NewRequestInfo() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare two RequestInfo structs
func equalRequestInfo(a, b *RequestInfo[any]) bool {
	return a.method == b.method &&
		a.url == b.url &&
		a.body == b.body &&
		a.authType == b.authType &&
		a.authCredentials.username == b.authCredentials.username &&
		a.authCredentials.password == b.authCredentials.password &&
		a.authCredentials.token == b.authCredentials.token &&
		equalHeaders(a.headers, b.headers)
}

// Helper function to compare two Headers maps
func equalHeaders(a, b Headers) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

// Test individual option functions
func TestWithAuthBasic(t *testing.T) {
	reqInfo := &RequestInfo[any]{}
	WithAuthBasic[any]("user", "pass")(reqInfo)
	if reqInfo.authType != AuthTypeBasic || reqInfo.authCredentials.username != "user" || reqInfo.authCredentials.password != "pass" {
		t.Errorf("WithAuthBasic() failed: %+v", reqInfo)
	}
}

func TestWithAuthBearer(t *testing.T) {
	reqInfo := &RequestInfo[any]{}
	WithAuthBearer[any]("some-token")(reqInfo)
	if reqInfo.authType != AuthTypeBearer || reqInfo.authCredentials.token != "some-token" {
		t.Errorf("WithAuthBearer() failed: %+v", reqInfo)
	}
}

func TestWithHeader(t *testing.T) {
	reqInfo := &RequestInfo[any]{headers: make(Headers)}
	WithHeader[any]("Custom-Header", "value")(reqInfo)
	if reqInfo.headers["Custom-Header"] != "value" {
		t.Errorf("WithHeader() failed: %+v", reqInfo.headers)
	}
}

func TestWithContentType(t *testing.T) {
	reqInfo := &RequestInfo[any]{headers: make(Headers)}
	WithContentType[any]("application/json")(reqInfo)
	if reqInfo.headers[HeaderContentType] != "application/json" {
		t.Errorf("WithContentType() failed: %+v", reqInfo.headers)
	}
}

func TestWithAccept(t *testing.T) {
	reqInfo := &RequestInfo[any]{headers: make(Headers)}
	WithAccept[any]("application/json")(reqInfo)
	if reqInfo.headers[HeaderAccept] != "application/json" {
		t.Errorf("WithAccept() failed: %+v", reqInfo.headers)
	}
}
