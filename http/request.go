package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// DoRequest Request Http and then convert to Response object
func DoRequest[T any](method string, url string, body any, bodySerializer func(body any) ([]byte, error), responseParser func(reader io.ReadCloser) (*T, error)) (*T, error) {
	// Serialize Request Body
	bodyReader, err := toBodyReader(body, bodySerializer)
	if err != nil {
		log.Println("parse request error:", err)
		return nil, err
	}

	// Make Http Request
	req, err := makeHttpRequest(method, url, bodyReader)
	if err != nil {
		log.Println("http request error:", err)
		return nil, err
	}

	// Send Http Request
	res, err := Client.Do(req)
	if err != nil {
		log.Println("client do error:", err)
		return nil, err
	}
	defer res.Body.Close()

	// Deserialize Response
	response, err := responseParser(res.Body)
	if err != nil {
		log.Println("parse response error:", err)
	}
	return response, err
}

func toBodyReader(body any, serializer func(body any) ([]byte, error)) (*bytes.Reader, error) {
	if body == nil {
		return nil, nil
	}

	requestBody, err := serializer(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(requestBody), nil
}

func makeHttpRequest(method string, url string, bodyReader *bytes.Reader) (*http.Request, error) {
	if bodyReader == nil {
		return http.NewRequest(method, url, nil)
	}
	return http.NewRequest(method, url, bodyReader)
}
