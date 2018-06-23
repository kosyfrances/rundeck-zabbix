package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Tests that getEnv func fatals when "TEST" is not set
	var value string
	key := "TEST"

	if value = os.Getenv(key); value != "" {
		os.Unsetenv(key)
		defer func() {
			os.Setenv(key, value)
		}()
	}

	if os.Getenv("TEST_GET_ENV") == "1" {
		getEnv(key)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestGetEnv")
	cmd.Env = append(os.Environ(), "TEST_GET_ENV=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)
	if ok && !e.Success() {
		return
	}

	t.Errorf("GetEnv() failed to exit if env %q is not present", key)
}

func TestGetZabbixKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"result": "fake zabbix key"}`)
	}))
	defer ts.Close()

	zabbixURL := ts.URL
	getZabbixKey(zabbixURL)

	if os.Getenv("ZABBIX_API_KEY") != "fake zabbix key" {
		t.Error("Expected Zabbix API key to be set in environment")
	}
}
