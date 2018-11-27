package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ZabbixAPI struct {
	URL      string
	Key      string
	User     string
	Password string
}

// MakeAPIRequest makes an API call to Zabbix.
// It returns a response object and an error object.
func (z *ZabbixAPI) MakeAPIRequest(payload map[string]interface{}) (resp *http.Response, err error) {
	// Build the request
	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Cannot marshal payload. Error: %v\n", err)
	}

	body := bytes.NewReader(b)
	req, err := http.NewRequest("GET", z.URL, body)
	if err != nil {
		log.Fatalf("Cannot make HTTP request. Error: %v\n", err)
	}

	req.Header.Set("Content-type", "application/json")
	// Send the request via a client
	resp, err = http.DefaultClient.Do(req)
	return resp, err
}

// SetAPIKey sets Zabbix Key in the ZabbixAPI struct if it is not already set
func (z *ZabbixAPI) SetAPIKey() error {
	if z.Key != "" {
		return nil
	}
	key, err := z.getAPIKey()
	if err != nil {
		return err
	}
	z.Key = key
	return nil
}

func (z *ZabbixAPI) getAPIKey() (string, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "user.login",
		"params": map[string]interface{}{
			"user":     z.User,
			"password": z.Password,
		},
		"id":   1,
		"auth": nil,
	}

	type Result struct {
		Key string `json:"result"`
	}

	var r Result

	resp, err := z.MakeAPIRequest(payload)
	if err != nil {
		return "", fmt.Errorf("cannot make Zabbix API call. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", fmt.Errorf("cannot decode response. Error: %v", err)
	}
	return r.Key, nil
}
