package zabbix

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test that we can get Zabbix API key from an API call,
// set it in ZabbixAPI struct
func TestGetAndSetZabbixKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"result": "fake_zabbix_key"}`)
	}))
	defer ts.Close()

	var z ZabbixAPI
	z.User = "ZABBIX_USER"
	z.Password = "ZABBIX_PASSWORD"
	z.URL = ts.URL

	// Get API Key
	key, err := z.getAPIKey()
	if err != nil {
		t.Fatalf("Process ran with err %v, want ZABBIX_API_KEY to be fake_zabbix_key", err)
		return
	}
	if key != "fake_zabbix_key" {
		t.Error("Expected key to be fake_zabbix_key")
	}

	// Set API Key
	z.SetAPIKey()

	if z.Key != "fake_zabbix_key" {
		t.Error("Expected ZabbixAPI.Key to be fake_zabbix_key")
	}
}

// Test that given a payload,
// we receive a response when we make an API call
func TestMakeZabbixAPIRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "apiinfo.version",
		"params":  map[string]string{},
		"id":      1,
	}

	var z ZabbixAPI
	z.URL = ts.URL

	resp, err := z.MakeAPIRequest(payload)
	if err != nil {
		t.Fatalf("Process ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}
}
