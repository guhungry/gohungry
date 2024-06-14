package json

import (
	"encoding/json"
	"gohungry/http"
	"io"
)

func Get[T any](url string) (*T, error) {
	return requestJson[T](http.MethodGet, url, nil)
}

func Post[T any](url string, body any) (*T, error) {
	return requestJson[T](http.MethodPost, url, body)
}

func requestJson[T any](method string, url string, body any) (*T, error) {
	var result T
	return http.DoRequest(method, url, body, json.Marshal, toResponseObject(&result))
}

func toResponseObject[T any](result *T) func(reader io.ReadCloser) (*T, error) {
	return func(reader io.ReadCloser) (*T, error) {
		if err := json.NewDecoder(reader).Decode(result); err != nil {
			return nil, err
		}
		return result, nil
	}
}
