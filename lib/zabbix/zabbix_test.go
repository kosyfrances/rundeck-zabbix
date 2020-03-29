package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const timeout = 1 * time.Second

// Test that we can get Zabbix API key from an API call,
// set it in API struct
func TestGetKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"result": "fake_zabbix_key"}`)
	}))
	defer ts.Close()

	api := API{
		URL:      ts.URL,
		User:     "ZABBIX_USER",
		Password: "ZABBIX_PASSWORD",
	}

	// Get API Key
	key, err := api.GetKey(timeout)
	if err != nil {
		t.Fatalf("Process ran with err %v, want ZABBIX_API_KEY to be fake_zabbix_key", err)
	}
	if key != "fake_zabbix_key" {
		t.Error("Expected key to be fake_zabbix_key")
	}
}

// Tests that we can get Resource info from a Zabbix API call
func TestGetHostsInfo(t *testing.T) {
	expected := HostResults{
		{
			HostID:      "10261",
			Host:        "dummy-host",
			Name:        "dummy-host",
			Description: "",
		},
	}
	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := struct {
			HostResults `json:"result"`
		}{expected}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	a := API{
		URL: ts.URL,
		Key: "fake_zabbix_key",
	}

	res, err := a.GetHostsInfo(timeout)
	if err != nil {
		t.Fatalf("process ran with err %v, want result to be %v", err, expected)
	}

	if !(reflect.DeepEqual(res, expected)) {
		t.Errorf("expected result to be %v but got %v", expected, res)
	}
}

func TestGetTriggersInfo(t *testing.T) {
	expected := TriggerResults{
		{
			Description: "Random trigger description",
			Hosts: []TriggerHost{
				TriggerHost{
					Name: "dummy-host",
				},
			},
		},
	}

	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := struct {
			TriggerResults `json:"result"`
		}{expected}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	a := API{
		URL: ts.URL,
		Key: "fake_zabbix_key",
	}

	res, err := a.GetTriggersInfo(timeout)
	if err != nil {
		t.Fatalf("process ran with err %v, want result to be %v", err, expected)
	}

	if !(reflect.DeepEqual(res, expected)) {
		t.Errorf("expected result to be %v but got %v", expected, res)
	}
}

func TestCreateClientUsingAPIKey(t *testing.T) {
	// Ensure that client cannot be created without both value and key
	type test struct {
		URL string
		key string
		err string
	}

	tests := []test{
		{
			URL: "fake_url",
			key: "",
			err: "zabbix API key missing. please run setup again",
		},
		{
			URL: "",
			key: "fake_key",
			err: "zabbix URL missing. please run setup again",
		},
		{
			URL: "",
			key: "",
			err: "zabbix URL missing. please run setup again",
		},
	}

	for _, tc := range tests {
		_, err := CreateClientUsingAPIKey(tc.URL, tc.key)
		if err == nil {
			t.Errorf("expected error to be %v but got %v", tc.err, err)
		}
	}

	// Ensure that client is can be created with value and key
	expected := &API{
		URL: "fake_url",
		Key: "fake_key",
	}

	a, err := CreateClientUsingAPIKey(expected.URL, expected.Key)
	if err != nil {
		t.Errorf("expected client to be %v but got %v", expected, a)
	}
}

func TestCreateClientUsingAuth(t *testing.T) {
	// Ensure that client is not created without incomplete values
	type test struct {
		URL      string
		user     string
		password string
	}

	tests := []test{
		{
			URL:      "",
			user:     "fake_user",
			password: "fake_password",
		},
		{
			URL:      "fake_url",
			user:     "",
			password: "fake_password",
		},
		{
			URL:      "fake_url",
			user:     "fake_user",
			password: "",
		},
	}

	for _, tc := range tests {
		_, err := CreateClientUsingAuth(tc.URL, tc.user, tc.password)
		if err == nil {
			t.Errorf("expected an error but got nil")
		}
	}

	// Ensure that client is created with required values
	expected := &API{
		URL:      "fake_url",
		User:     "fake_user",
		Password: "fake_password",
	}

	a, err := CreateClientUsingAuth(expected.URL, expected.User, expected.Password)
	if err != nil {
		t.Errorf("expected client to be %v but got %v", expected, a)
	}
}

func TestAcknowledgeEvent(t *testing.T) {
	expected := map[string][]int{"eventids": []int{49}}

	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := struct {
			Event map[string][]int `json:"result"`
		}{expected}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	a := API{
		URL: ts.URL,
		Key: "fake_zabbix_key",
	}

	res, err := a.AcknowledgeEvent("49", "fake message", timeout)
	if err != nil {
		t.Fatalf("process ran with err %v, want result to be %v", err, []int{49})
	}

	if !reflect.DeepEqual(res, []int{49}) {
		t.Errorf("expected Zabbix acknowledgement event ID to be %v", []int{49})
	}

}
