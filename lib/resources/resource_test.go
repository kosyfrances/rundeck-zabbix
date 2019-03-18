package resources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

// Test that we can get Rundeck resource model document in Yaml format
// with hosts' information from Zabbix
func TestMakeResource(t *testing.T) {
	expected := zabbix.HostResults{
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
			zabbix.HostResults `json:"result"`
		}{expected}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			t.Fatalf("unable to write response. error: %v", err)
		}
	}))
	defer ts.Close()

	// generate temp file
	tmpfile, err := ioutil.TempFile("", "testfile.*.yml")
	if err != nil {
		t.Fatalf("unable to create temporary file. error: %v", err)
	}

	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()           // close file

	a, err := zabbix.CreateClientUsingAPIKey(ts.URL, "fake_zabbix_key")
	if err != nil {
		t.Fatalf("process ran with err %v", err)
	}

	res, err := a.GetHostsInfo()
	if err != nil {
		t.Fatalf("cannot get hosts info.\n%v", err)
	}

	Make(res, tmpfile.Name())
	if err != nil {
		t.Fatalf("expected resource model document to be generated. error: %v", err)
	}

	// check the file's content
	e := `
dummy-host:
  hostname: dummy-host
  nodename: dummy-host
  description: ""
  osArch: ""
  osFamily: ""
  osName: ""
  osVersion: ""
  tags: ""
  username: ""
  ssh-keypath: ""`

	c, err := ioutil.ReadFile(tmpfile.Name())

	if err != nil {
		t.Fatalf("expected file content to be %v\n\n error: %v", e, err)
	}

	content := strings.Trim(string(c), "\n")
	expectedContent := strings.Trim(e, "\n")

	if content != expectedContent {
		t.Fatalf("test returned \n%v\n expected file content to be \n%v", content, expectedContent)
	}
}
