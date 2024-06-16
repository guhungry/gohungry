package xml

import (
	"encoding/xml"
	"gohungry/http"
	"io"
)

// contentType specifies the MIME type for XML content.
const contentType = "application/xml"

// Get performs an HTTP GET, decodes XML response into 'Response'.
// 'url' is the request target, 'options' customize the request.
func Get[Response any](url string, options ...http.RequestInfoOption[Response]) (*Response, error) {
	return requestXML[Response](http.MethodGet, url, nil, options...)
}

// Post performs an HTTP POST with XML body, decodes response into 'Response'.
// 'url' is the request target, 'body' is the payload, 'options' customize the request.
func Post[Response any](url string, body any, options ...http.RequestInfoOption[Response]) (*Response, error) {
	return requestXML[Response](http.MethodPost, url, body, options...)
}

// requestXML sends an HTTP request and decodes XML response into 'Response'.
// 'method' is the HTTP method, 'url' is the request target, 'body' is the payload for POST,
// 'options' customize the request.
func requestXML[Response any](method string, url string, body any, options ...http.RequestInfoOption[Response]) (*Response, error) {
	options = append(options,
		http.WithAccept[Response](contentType),
		http.WithContentType[Response](contentType),
	)
	request := http.NewRequestInfo(method, url, body, xml.Marshal, toResponseObject[Response](), options...)
	return http.DoRequest(request)
}

// toResponseObject returns a function that decodes XML into 'Response'.
func toResponseObject[Response any]() http.ResponseBodyParser[Response] {
	return func(reader io.ReadCloser) (*Response, error) {
		var result Response
		decoder := xml.NewDecoder(reader)
		if err := decoder.Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}
