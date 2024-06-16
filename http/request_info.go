package http

import "io"

type RequestBodySerializer func(body any) ([]byte, error)
type ResponseBodyParser[T any] func(reader io.ReadCloser) (*T, error)

// RequestInfo contains HTTP request information including methods for serializing
// the request body and deserializing the response body, as well as HTTP headers.
// Use NewRequestInfo to initialize this struct.
type RequestInfo[T any] struct {
	method         string
	url            string
	body           any
	bodySerializer RequestBodySerializer
	responseParser ResponseBodyParser[T]
	authType       string
	authUsername   string
	authPassword   string
	authToken      string
	headers        map[string]string // HTTP Headers
}

// RequestInfoOption is a function type that modifies RequestInfo.
type RequestInfoOption[T any] func(c *RequestInfo[T])

// NewRequestInfo initializes and returns a new RequestInfo with the given parameters.
func NewRequestInfo[T any](method string, url string, body any, bodySerializer RequestBodySerializer, responseParser ResponseBodyParser[T], options ...RequestInfoOption[T]) *RequestInfo[T] {
	result := &RequestInfo[T]{
		method:         method,
		url:            url,
		body:           body,
		bodySerializer: bodySerializer,
		responseParser: responseParser,
		headers:        make(map[string]string),
	}

	for _, option := range options {
		option(result)
	}
	return result
}

// WithAuthBasic adds basic authentication to the HTTP request.
func WithAuthBasic[T any](username, password string) RequestInfoOption[T] {
	return func(c *RequestInfo[T]) {
		c.authType = AuthTypeBasic
		c.authUsername = username
		c.authPassword = password
	}
}

// WithAuthBearer adds bearer token authentication to the HTTP request.
func WithAuthBearer[T any](token string) RequestInfoOption[T] {
	return func(c *RequestInfo[T]) {
		c.authType = AuthTypeBearer
		c.authToken = token
	}
}

// WithHeader adds a key-value pair to the HTTP headers of the request.
func WithHeader[T any](key, value string) RequestInfoOption[T] {
	return func(c *RequestInfo[T]) {
		c.headers[key] = value
	}
}

// WithContentType adds a Content-Type header to the HTTP request.
func WithContentType[T any](value string) RequestInfoOption[T] {
	return WithHeader[T](HeaderContentType, value)
}

// WithAccept adds an Accept header to the HTTP request.
func WithAccept[T any](value string) RequestInfoOption[T] {
	return WithHeader[T](HeaderAccept, value)
}
