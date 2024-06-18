package http

import (
	"bytes"
	"log"
	"net/http"
)

// DoRequest executes an HTTP request and decodes the response into 'Request'.
// 'data' contains request configurations and handlers.
func DoRequest[Request any](data *RequestInfo[Request]) (*Request, error) {
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

// setHeaders applies provided headers to the HTTP request.
func setHeaders(req *http.Request, headers Headers) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

// setAuth configures authentication for the HTTP request using 'data'.
func setAuth[Request any](req *http.Request, data *RequestInfo[Request]) {
	auth := data.authCredentials
	switch data.authType {
	case AuthTypeBasic:
		req.SetBasicAuth(auth.username, auth.password)
	case AuthTypeBearer:
		req.Header.Add(HeaderAuthorization, string(AuthTypeBearer)+" "+auth.token)
	default:
		// No authentication required
	}
}

// toBodyReader creates a reader for serialized request body.
// 'body' is the payload, 'serializer' converts it to a byte slice.
func toBodyReader(body any, serializer RequestBodySerializer) (*bytes.Reader, error) {
	if body == nil {
		return bytes.NewReader([]byte{}), nil
	}

	requestBody, err := serializer(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(requestBody), nil
}
