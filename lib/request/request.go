package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// ZabbixHeaderKey is the header key for Zabbix API request
	// The header for Zabbix is "Content-type: application/json"
	ZabbixHeaderKey = "Content-type"
	// RundeckHeaderKey is the header key for Rundeck API request
	// The header for Rundeck is "Accept: application/json"
	RundeckHeaderKey = "Accept"
)

func build(method, URL string, payload interface{}) (*http.Request, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload. error: %v", err)
	}

	body := bytes.NewReader(b)
	return http.NewRequest(method, URL, body)
}

// Make "makes" an API request.
// It returns a response object and an error object.
func Make(headerKey, method, URL string, timeout time.Duration, payload interface{}) (*http.Response, error) {
	// Build the request
	req, err := build(method, URL, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot create HTTP request. error: %v", err)
	}

	req.Header.Set(headerKey, "application/json")

	// Send the request via a client
	client := http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}
