package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func buildRequest(method, URL string, payload interface{}) (*http.Request, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload. error: %v", err)
	}

	body := bytes.NewReader(b)
	return http.NewRequest(method, URL, body)
}

// MakeRundeckRequest makes an API request.
// It sets the header "Accept" as "application/json".
// It returns a response object and an error object.
func MakeRundeckRequest(method, URL string, payload interface{}) (*http.Response, error) {
	// Build the request
	req, err := buildRequest(method, URL, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot create HTTP request. error: %v", err)
	}

	req.Header.Set("Accept", "application/json")

	// Send the request via a client
	return http.DefaultClient.Do(req)
}

// MakeZabbixRequest makes an API request.
// It sets the header "Content-type" as "application/json".
// It returns a response object and an error object.
func MakeZabbixRequest(method, URL string, payload interface{}) (*http.Response, error) {
	// Build the request
	req, err := buildRequest(method, URL, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot create HTTP request. error: %v", err)
	}

	req.Header.Set("Content-type", "application/json")

	// Send the request via a client
	return http.DefaultClient.Do(req)
}

// DumpToFile dumps data, a list of bytes to file given.
// If the filePath doesn't exist, it creates it, or appends to the file
func DumpToFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file. %v", err)
	}

	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("cannot write to file. %v", err)
	}
	return nil
}
