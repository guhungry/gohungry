package json

import (
	"encoding/json"
	"gohungry/http"
	"io"
)

// contentType is the MIME type for JSON content.
const contentType = "application/json"

// Get sends an HTTP GET request and decodes the JSON response into type T.
func Get[T any](url string, options ...http.RequestInfoOption[T]) (*T, error) {
	return requestJSON[T](http.MethodGet, url, nil, options...)
}

// Post sends an HTTP POST request with JSON body and decodes the response into type T.
func Post[T any](url string, body any, options ...http.RequestInfoOption[T]) (*T, error) {
	return requestJSON[T](http.MethodPost, url, body, options...)
}

// requestJSON sets up and sends an HTTP request with JSON content and decodes the response.
func requestJSON[T any](method string, url string, body any, options ...http.RequestInfoOption[T]) (*T, error) {
	options = append(options,
		http.WithAccept[T](contentType),
		http.WithContentType[T](contentType),
	)
	request := http.NewRequestInfo(method, url, body, json.Marshal, toResponseObject[T](), options...)
	return http.DoRequest[T](request)
}

// toResponseObject creates a function that decodes a JSON response into type T.
func toResponseObject[T any]() http.ResponseBodyParser[T] {
	return func(reader io.ReadCloser) (*T, error) {
		var result T
		if err := json.NewDecoder(reader).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}
