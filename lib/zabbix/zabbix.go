package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kosyfrances/rundeck-zabbix/lib/utils"
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
	Auth    string      `json:"auth,omitempty"`
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

// API response Error struct
type apiError struct {
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

// CreateClientUsingAPIKey creates new API client using key
func CreateClientUsingAPIKey(URL, key string) (*API, error) {
	if URL == "" {
		return nil, fmt.Errorf("zabbix URL missing. please run setup again")
	}

	if key == "" {
		return nil, fmt.Errorf("zabbix API key missing. please run setup again")
	}

	return &API{
		URL: URL,
		Key: key,
	}, nil
}

// CreateClientUsingAuth creates new API client using username and password
func CreateClientUsingAuth(URL, user, password string) (*API, error) {
	if URL == "" {
		return nil, fmt.Errorf("zabbix URL not found. run setup again")
	}

	// Ensure that either key, username and password are not empty at the same time
	if user == "" || password == "" {
		return nil, fmt.Errorf("zabbix username and password must be present. run setup again")
	}

	return &API{
		URL:      URL,
		User:     user,
		Password: password,
	}, nil
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

// GetKey gets Zabbix API key
func (a *API) GetKey() (string, error) {
	params := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{
		User:     a.User,
		Password: a.Password,
	}

	payload := a.BuildPayload(params, "user.login")

	var r struct {
		Key string   `json:"result"`
		Err apiError `json:"error"`
	}

	resp, err := utils.MakeRequest(http.MethodGet, a.URL, payload)

	if err != nil {
		return "", fmt.Errorf("cannot make Zabbix API call. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", fmt.Errorf("cannot decode response. Error: %v", err)
	}
	//  Check if the response contains an error
	if r.Err != (apiError{}) {
		return "", fmt.Errorf("%v %v", r.Err.Message, r.Err.Data)
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

	resp, err := utils.MakeRequest(http.MethodGet, a.URL, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot make API request. error: %v", err)
	}

	var r struct {
		Results `json:"result"`
		Err     apiError `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response. error: %v", err)
	}

	//  Check if the response contains an error
	if r.Err != (apiError{}) {
		return r.Results, fmt.Errorf("%v %v", r.Err.Message, r.Err.Data)
	}
	return r.Results, nil
}
