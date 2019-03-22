package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test that given a payload,
// we receive a response when we make an API call
func TestMakeRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	// Get request
	resp, err := MakeRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Fatalf("Process ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}

	// Post request
	resp, err = MakeRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Fatalf("Process ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Error("Expected response status code to be 200")
	}

}
