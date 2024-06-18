package http

import "io"

// RequestBodySerializer serializes a request body into a byte slice.
type RequestBodySerializer func(body any) ([]byte, error)

// ResponseBodyParser decodes a response body into a specified type 'Response'.
type ResponseBodyParser[Response any] func(reader io.ReadCloser) (*Response, error)

// RequestInfo holds configurations for an HTTP request, including serialization
// for the request body and deserialization for the response body.
// Use NewRequestInfo to initialize this struct.
type RequestInfo[Response any] struct {
	method          string
	url             string
	body            any
	bodySerializer  RequestBodySerializer
	responseParser  ResponseBodyParser[Response]
	authType        AuthType
	authCredentials AuthCredentials
	headers         Headers // HTTP Headers
}

// AuthCredentials holds authentication credentials.
type AuthCredentials struct {
	username string
	password string
	token    string
}

// Headers represents HTTP headers as a map.
type Headers map[string]string

// RequestInfoOption modifies a RequestInfo instance.
type RequestInfoOption[Response any] func(c *RequestInfo[Response])

// NewRequestInfo creates a new RequestInfo with specified parameters and options.
func NewRequestInfo[Response any](method string, url string, body any, bodySerializer RequestBodySerializer, responseParser ResponseBodyParser[Response], options ...RequestInfoOption[Response]) *RequestInfo[Response] {
	result := &RequestInfo[Response]{
		method:         method,
		url:            url,
		body:           body,
		bodySerializer: bodySerializer,
		responseParser: responseParser,
		headers:        make(Headers),
	}

	for _, option := range options {
		option(result)
	}
	return result
}

// WithAuthBasic sets basic authentication credentials for the request.
func WithAuthBasic[Response any](username, password string) RequestInfoOption[Response] {
	return func(c *RequestInfo[Response]) {
		c.authType = AuthTypeBasic
		c.authCredentials = AuthCredentials{username: username, password: password}
	}
}

// WithAuthBearer sets bearer token authentication for the request.
func WithAuthBearer[Response any](token string) RequestInfoOption[Response] {
	return func(c *RequestInfo[Response]) {
		c.authType = AuthTypeBearer
		c.authCredentials = AuthCredentials{token: token}
	}
}

// WithHeader adds a header to the request.
func WithHeader[Response any](key, value string) RequestInfoOption[Response] {
	return func(c *RequestInfo[Response]) {
		c.headers[key] = value
	}
}

// WithContentType sets the Content-Type header for the request.
func WithContentType[Response any](value string) RequestInfoOption[Response] {
	return WithHeader[Response](HeaderContentType, value)
}

// WithAccept sets the Accept header for the request.
func WithAccept[Response any](value string) RequestInfoOption[Response] {
	return WithHeader[Response](HeaderAccept, value)
}
