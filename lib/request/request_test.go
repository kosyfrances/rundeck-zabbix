package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const timeout = 1 * time.Second

// Test that given a payload,
// we receive a response when we make an API call
func TestMakeRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	// Get Rundeck request
	resp, err := Make(RundeckHeaderKey, http.MethodGet, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeRundeckRequest with Get method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}

	// Get Zabbix request
	resp, err = Make(ZabbixHeaderKey, http.MethodGet, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeZabbixRequest with Get method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}

	// Post Rundeck request
	resp, err = Make(RundeckHeaderKey, http.MethodPost, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeRundeckRequest with Post method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}

	// Post Zabbix request
	resp, err = Make(ZabbixHeaderKey, http.MethodPost, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeZabbixRequest Post method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}
}

// Test that given a delay,
// we do not receive a timeout earlier when we make an API call
func TestMakeRequestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("actual"))
	}))
	defer ts.Close()

	// Get Rundeck request
	resp, err := Make(RundeckHeaderKey, http.MethodGet, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeRundeckRequest with Get method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}
}
