package http

import (
	"bytes"
	"log"
	"net/http"
)

// DoRequest sends an HTTP request and converts the response into a Response object of type T.
func DoRequest[T any](data *RequestInfo[T]) (*T, error) {
	// Serialize Request Body
	bodyReader, err := toBodyReader(data.body, data.bodySerializer)
	if err != nil {
		log.Println("parse request error:", err)
		return nil, err
	}

	// Make Http Request
	req, err := http.NewRequest(data.method, data.url, bodyReader)
	if err != nil {
		log.Println("http request error:", err)
		return nil, err
	}
	setAuth(req, data)
	setHeaders(req, data.headers)

	// Send Http Request
	res, err := Client.Do(req)
	if err != nil {
		log.Println("client do error:", err)
		return nil, err
	}
	defer res.Body.Close()

	// Deserialize Response
	response, err := data.responseParser(res.Body)
	if err != nil {
		log.Println("parse response error:", err)
		return nil, err
	}
	return response, nil
}

// setHeaders sets the headers for the HTTP request.
func setHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

// setAuth sets the authentication for the HTTP request based on the provided RequestInfo.
func setAuth[T any](req *http.Request, data *RequestInfo[T]) {
	switch data.authType {
	case AuthTypeBasic:
		req.SetBasicAuth(data.authUsername, data.authPassword)
	case AuthTypeBearer:
		req.Header.Add(HeaderAuthorization, data.authType+" "+data.authToken)
	default:
		// No authentication required
	}
}

// toBodyReader serializes the body using the provided serializer function and returns a bytes.Reader.
func toBodyReader(body any, serializer func(body any) ([]byte, error)) (*bytes.Reader, error) {
	if body == nil {
		return bytes.NewReader([]byte{}), nil
	}

	requestBody, err := serializer(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(requestBody), nil
}
