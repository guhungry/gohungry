package xml

import (
	"encoding/xml"
	"gohungry/http"
	"io"
)

func Get[T any](url string) (*T, error) {
	return requestXml[T](http.MethodGet, url, nil)
}

func Post[T any](url string, body any) (*T, error) {
	return requestXml[T](http.MethodPost, url, body)
}

func requestXml[T any](method string, url string, body any) (*T, error) {
	var result T
	return http.DoRequest(method, url, body, xml.Marshal, toResponseObject(&result))
}

func toResponseObject[T any](result *T) func(reader io.ReadCloser) (*T, error) {
	return func(reader io.ReadCloser) (*T, error) {
		if err := xml.NewDecoder(reader).Decode(result); err != nil {
			return nil, err
		}
		return result, nil
	}
}
