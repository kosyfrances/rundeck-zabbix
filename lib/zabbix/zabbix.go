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

// hostResult struct holds a Zabbix host info
type hostResult struct {
	HostID      string `json:"hostid"`
	Host        string `json:"host"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// HostResults struct holds a list of Zabbix hosts
type HostResults []hostResult

// TriggerHost struct holds info about a Zabbix Host in a Trigger
type TriggerHost struct {
	Name string `json:"name"`
}

// triggerResult struct holds a Zabbix trigger info
type triggerResult struct {
	Description string        `json:"description"`
	Hosts       []TriggerHost `json:"hosts"`
}

// TriggerResults struct holds a list of Zabbix triggers
type TriggerResults []triggerResult

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
		Key string    `json:"result"`
		Err *apiError `json:"error"`
	}

	resp, err := utils.MakeZabbixRequest(http.MethodGet, a.URL, payload)

	if err != nil {
		return "", fmt.Errorf("cannot make Zabbix API call. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", fmt.Errorf("cannot decode response. Error: %v", err)
	}
	//  Check if the response contains an error
	if r.Err != nil {
		return "", fmt.Errorf("%v %v", r.Err.Message, r.Err.Data)
	}

	return r.Key, nil
}

// GetHostsInfo gets hosts information from Zabbix
func (a *API) GetHostsInfo() (HostResults, error) {
	params := struct {
		Output []string `json:"output"`
	}{
		Output: []string{"host", "name", "description"},
	}

	payload := a.BuildPayload(params, "host.get")

	resp, err := utils.MakeZabbixRequest(http.MethodGet, a.URL, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot make API request. error: %v", err)
	}

	var r struct {
		HostResults `json:"result"`
		Err         *apiError `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response. error: %v", err)
	}

	//  Check if the response contains an error
	if r.Err != nil {
		return r.HostResults, fmt.Errorf("%v %v", r.Err.Message, r.Err.Data)
	}
	return r.HostResults, nil
}

// GetTriggersInfo gets triggers information for hosts from Zabbix
func (a *API) GetTriggersInfo() (TriggerResults, error) {
	params := struct {
		Output            []string `json:"output"`
		SelectHosts       []string `json:"selectHosts"`
		Active            string   `json:"active"`
		ExpandDescription string   `json:"expandDescription"`
	}{
		Output:      []string{"description"},
		SelectHosts: []string{"name"},
	}

	payload := a.BuildPayload(params, "trigger.get")

	resp, err := utils.MakeZabbixRequest(http.MethodGet, a.URL, payload)

	if err != nil {
		return nil, fmt.Errorf("cannot make API request. error: %v", err)
	}

	var r struct {
		TriggerResults `json:"result"`
		Err            *apiError `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response. error: %v", err)
	}

	//  Check if the response contains an error
	if r.Err != nil {
		return r.TriggerResults, fmt.Errorf("%v %v", r.Err.Message, r.Err.Data)
	}
	return r.TriggerResults, nil
}

// AcknowledgeEvent acknowledges a single event and leaves a message.
// It returns the event ID and an error if any.
func (a *API) AcknowledgeEvent(eventID, message string) ([]int, error) {
	params := struct {
		EventIDs string `json:"eventids"`
		// Event update action(s). This is bitmask field, any combination of values are acceptable.
		// Possible values:
		// 1 - close problem;
		// 2 - acknowledge event;
		// 4 - add message;
		// 8 - change severity.
		Action  int    `json:"action"`
		Message string `json:"message"`
	}{
		EventIDs: eventID,
		// 6 == 4 (add message) + 2 (ack event)
		Action:  6,
		Message: message,
	}

	payload := a.BuildPayload(params, "event.acknowledge")

	resp, err := utils.MakeZabbixRequest(http.MethodPost, a.URL, payload)

	if err != nil {
		return nil, fmt.Errorf("cannot make API request. error: %v", err)
	}

	var r struct {
		Result struct {
			EventIDs []int `json:"eventids"`
		} `json:"result"`
		Err *apiError `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response from Zabbix service. error: %v", err)
	}

	//  Check if the response contains an error
	if r.Err != nil {
		return nil, fmt.Errorf("Zabbix service returned an error response. %v %v", r.Err.Message, r.Err.Data)
	}

	return r.Result.EventIDs, nil
}
