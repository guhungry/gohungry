package xml

import (
	"encoding/xml"
	"gohungry/http"
	"io"
)

// contentType is the MIME type for XML content.
const contentType = "application/xml"

// Get sends an HTTP GET request and decodes the XML response into type T.
func Get[T any](url string, options ...http.RequestInfoOption[T]) (*T, error) {
	return requestXML[T](http.MethodGet, url, nil, options...)
}

// Post sends an HTTP POST request with XML body and decodes the response into type T.
func Post[T any](url string, body any, options ...http.RequestInfoOption[T]) (*T, error) {
	return requestXML[T](http.MethodPost, url, body, options...)
}

// requestXML sets up and sends an HTTP request with XML content and decodes the response.
func requestXML[T any](method string, url string, body any, options ...http.RequestInfoOption[T]) (*T, error) {
	options = append(options,
		http.WithAccept[T](contentType),
		http.WithContentType[T](contentType),
	)
	request := http.NewRequestInfo(method, url, body, xml.Marshal, toResponseObject[T](), options...)
	return http.DoRequest(request)
}

// toResponseObject creates a function that decodes an XML response into type T.
func toResponseObject[T any]() http.ResponseBodyParser[T] {
	return func(reader io.ReadCloser) (*T, error) {
		var result T
		if err := xml.NewDecoder(reader).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}
