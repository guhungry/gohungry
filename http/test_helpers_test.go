package http

import (
	"encoding/json"
	"io"
)

// Dummy implementations for testing.
func dummyRequestBodySerializer(body any) ([]byte, error) {
	return []byte("dummy body"), nil
}

func dummyResponseBodyParser[Response any](reader io.ReadCloser) (*Response, error) {
	var result Response
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
