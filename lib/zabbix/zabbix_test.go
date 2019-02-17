package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Test that we can get Zabbix API key from an API call,
// set it in API struct
func TestGetAndSetKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"result": "fake_zabbix_key"}`)
	}))
	defer ts.Close()

	var api API
	api.User = "ZABBIX_USER"
	api.Password = "ZABBIX_PASSWORD"
	api.URL = ts.URL

	// Get API Key
	key, err := api.getKey()
	if err != nil {
		t.Fatalf("Process ran with err %v, want ZABBIX_API_KEY to be fake_zabbix_key", err)
	}
	if key != "fake_zabbix_key" {
		t.Error("Expected key to be fake_zabbix_key")
	}

	// Set API Key
	api.SetKey()

	if api.Key != "fake_zabbix_key" {
		t.Error("Expected ZabbixAPI.Key to be fake_zabbix_key")
	}
}

// Test that given a payload,
// we receive a response when we make an API call
func TestMakeRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	var api API
	api.URL = ts.URL

	payload := api.BuildPayload(nil, "apiinfo.version")
	resp, err := api.MakeRequest(payload)
	if err != nil {
		t.Fatalf("Process ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}
}

// Tests that we can get Resource info from a Zabbix API call
func TestGetHostsInfo(t *testing.T) {
	expected := Results{
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
			Results `json:"result"`
		}{expected}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	a := NewAPI(ts.URL, "fake_zabbix_key")
	res, err := a.GetHostsInfo()
	if err != nil {
		t.Fatalf("process ran with err %v, want result to be %v", err, expected)
	}

	if !(reflect.DeepEqual(res, expected)) {
		t.Errorf("expected result to be %v but got %v", expected, res)
	}
}
