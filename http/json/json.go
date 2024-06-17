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
	request := http.NewRequestInfo(method, url, body, json.Marshal, toResponseObject[Response](), options...)
	return http.DoRequest[Response](request)
}

// toResponseObject returns a function that decodes JSON into 'Response'.
func toResponseObject[Response any]() http.ResponseBodyParser[Response] {
	return func(reader io.ReadCloser) (*Response, error) {
		var result Response
		if err := json.NewDecoder(reader).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}
