// Package json provides utilities for making HTTP requests with JSON payloads
// and parsing JSON responses. It includes functions for performing GET and POST
// requests with JSON content, along with customizable options for request configuration.
package json

import (
	"encoding/json"
	"github.com/guhungry/gohungry/http"
	"io"
)

// contentType specifies the MIME type for JSON content.
const contentType = "application/json"

// Get performs an HTTP GET, decodes JSON response into 'Response'.
// 'url' is the request target, 'options' customize the request.
func Get[Response any](url string, options ...http.RequestInfoOption[Response]) (*Response, error) {
	return requestJSON[Response](http.MethodGet, url, nil, options...)
}

// Post performs an HTTP POST with JSON body, decodes response into 'Response'.
// 'url' is the request target, 'body' is the payload, 'options' customize the request.
func Post[Response any](url string, body any, options ...http.RequestInfoOption[Response]) (*Response, error) {
	return requestJSON[Response](http.MethodPost, url, body, options...)
}

// requestJSON sends an HTTP request and decodes JSON response into 'Response'.
// 'method' is the HTTP method, 'url' is the request target, 'body' is the payload for POST,
// 'options' customize the request.
func requestJSON[Response any](method string, url string, body any, options ...http.RequestInfoOption[Response]) (*Response, error) {
	options = append(options,
		http.WithAccept[Response](contentType),
		http.WithContentType[Response](contentType),
	)
	request := http.NewRequestInfo(method, url, body, json.Marshal, toResponseObject[Response], options...)
	return http.DoRequest[Response](request)
}

// toResponseObject decodes JSON into 'Response'.
func toResponseObject[Response any](reader io.ReadCloser) (*Response, error) {
	var result Response
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
