package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const timeout = 1 * time.Second

func TestBuildRundeckURLEndpoint(t *testing.T) {
	baseURL := "http://rundeckbaseurl:4440"
	endpoint := "/api/18/job"
	fullURL := "http://rundeckbaseurl:4440/api/18/job"

	URL, err := BuildRundeckURLEndpoint(baseURL, endpoint)
	if err != nil {
		t.Fatalf("Process ran with err %v, want URL to be %s", err, fullURL)
	}

	if URL != fullURL {
		t.Errorf("Expected key to be %s", fullURL)
	}
}

func TestGetRundeckJobID(t *testing.T) {
	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"id":"fake-job-id"}]`)
	}))
	defer ts.Close()

	ID, err := GetRundeckJobID(ts.URL, timeout)
	if err != nil {
		t.Fatalf("Process ran with err %v, want ID to be fake-job-id", err)
	}

	if ID != "fake-job-id" {
		t.Error("Expected ID to be fake-job-id")
	}
}

func TestExecuteRundeckJob(t *testing.T) {
	expected := JobExecResponse{
		ID:      4,
		Status:  "running",
		Project: "test-project",
		Job: job{
			Name: "fake-job",
		},
	}
	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expected)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	resp, err := ExecuteRundeckJob(ts.URL, timeout)
	if err != nil {
		t.Fatalf("Process ran with err %v", err)
	}
	if *resp != expected {
		t.Errorf("expected rundeck execution to return %v", expected)
	}
}
