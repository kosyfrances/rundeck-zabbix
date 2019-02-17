package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// API struct details for Zabbix
type API struct {
	URL      string
	Key      string
	User     string
	Password string
}

// Payload struct holds info needed for an API call
type Payload struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Auth    string      `json:"auth"`
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
}

// Result struct holds a Zabbix host info
type Result struct {
	HostID      string `json:"hostid"`
	Host        string `json:"host"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Results struct holds a list of Zabbix host info
type Results []Result

// NewAPI builds new API struct
func NewAPI(URL, key string) API {
	return API{
		URL: URL,
		Key: key,
	}
}

// BuildPayload builds payload with params and method given
func (a *API) BuildPayload(params interface{}, method string) Payload {
	return Payload{
		Method:  method,
		Params:  params,
		Auth:    a.Key,
		JSONRPC: "2.0",
		ID:      1,
	}
}

// MakeRequest makes an API call to Zabbix.
// It returns a response object and an error object.
func (a *API) MakeRequest(payload Payload) (*http.Response, error) {
	// Build the request
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload. error: %v", err)
	}

	body := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodGet, a.URL, body)
	if err != nil {
		return nil, fmt.Errorf("cannot make HTTP request. error: %v", err)
	}

	req.Header.Set("Content-type", "application/json")
	// Send the request via a client
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

// SetKey sets Zabbix Key in the API struct if it is not already set
func (a *API) SetKey() error {
	if a.Key != "" {
		return nil
	}
	key, err := a.getKey()
	if err != nil {
		return err
	}
	a.Key = key
	return nil
}

func (a *API) getKey() (string, error) {
	params := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{
		User:     a.User,
		Password: a.Password,
	}

	payload := a.BuildPayload(params, "user.login")

	var r struct {
		Key string `json:"result"`
	}

	resp, err := a.MakeRequest(payload)
	if err != nil {
		return "", fmt.Errorf("cannot make Zabbix API call. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", fmt.Errorf("cannot decode response. Error: %v", err)
	}
	return r.Key, nil
}

// GetHostsInfo gets hosts information from Zabbix
func (a *API) GetHostsInfo() (Results, error) {
	params := struct {
		Output         []string `json:"output"`
		SelectTriggers []string `json:"selectTriggers"`
	}{
		Output:         []string{"host", "name", "description"},
		SelectTriggers: []string{"description", "status"},
	}

	payload := a.BuildPayload(params, "host.get")

	resp, err := a.MakeRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot make API request. error: %v", err)
	}

	var r struct {
		Results `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response. error: %v", err)
	}

	return r.Results, nil
}
