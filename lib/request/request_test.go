package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const timeout = 1 * time.Second

// Test that we receive a response when we make an API call
func TestMakeRequest(t *testing.T) {
	type test struct {
		msg    string
		key    string
		method string
		want   int
	}

	tests := []test{
		{
			msg:    "MakeRundeckRequest",
			key:    RundeckHeaderKey,
			method: http.MethodGet,
			want:   200,
		},
		{
			msg:    "MakeRundeckRequest",
			key:    RundeckHeaderKey,
			method: http.MethodPost,
			want:   200,
		},
		{
			msg:    "MakeZabbixRequest",
			key:    ZabbixHeaderKey,
			method: http.MethodGet,
			want:   200,
		},
		{
			msg:    "MakeZabbixRequest",
			key:    ZabbixHeaderKey,
			method: http.MethodPost,
			want:   200,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	for _, tc := range tests {
		resp, err := Make(RundeckHeaderKey, http.MethodGet, ts.URL, timeout, nil)
		if err != nil {
			t.Fatalf("%s with %s method ran with err %v, want response", tc.msg, tc.method, err)
		}

		if resp.StatusCode != tc.want {
			t.Fatalf("%s with %s method returned status code %v, want %v", tc.msg, tc.method, resp.StatusCode, tc.want)
		}
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

	// Get Rundeck request does not timeout
	resp, err := Make(RundeckHeaderKey, http.MethodGet, ts.URL, timeout, nil)
	if err != nil {
		t.Fatalf("MakeRundeckRequest with Get method ran with err %v, want response", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Fatalf("MakeRundeckRequest with Get method returned status code %v, want %v", resp.StatusCode, 200)
	}

	// Get Rundeck request times out
	resp, err = Make(RundeckHeaderKey, http.MethodGet, ts.URL, 50*time.Millisecond, nil)
	if err == nil {
		t.Fatalf("MakeRundeckRequest with Get method ran without error, want error %v", err)
		return
	}

	if resp != nil {
		t.Fatalf("Expected MakeRundeckRequest with lower timeout to be nil")
	}
}
