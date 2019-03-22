package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MakeRequest makes an API request.
// It returns a response object and an error object.
func MakeRequest(method, URL string, payload interface{}) (*http.Response, error) {
	// Build the request
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload. error: %v", err)
	}

	body := bytes.NewReader(b)
	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, fmt.Errorf("cannot create HTTP request. error: %v", err)
	}

	req.Header.Set("Content-type", "application/json")
	// Send the request via a client
	return http.DefaultClient.Do(req)
}
