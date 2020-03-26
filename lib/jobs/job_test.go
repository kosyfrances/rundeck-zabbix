package jobs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

const timeout = 1 * time.Second

// Test that we can get Rundeck job model document in Yaml format
// with triggers' information from Zabbix
func TestMakeJob(t *testing.T) {
	expected := zabbix.TriggerResults{
		{
			Description: "random trigger on dummy-host",
			Hosts: []zabbix.TriggerHost{
				{Name: "dummy-host"},
			},
		}, {
			Description: "prefix random trigger on dummy-host2",
			Hosts: []zabbix.TriggerHost{
				{Name: "dummy-host2"},
			},
		},
	}

	// mock api call
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := struct {
			zabbix.TriggerResults `json:"result"`
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

	res, err := a.GetTriggersInfo(timeout)
	if err != nil {
		t.Fatalf("cannot get triggers info.\n%v", err)
	}

	// Without prefix for triggers
	err = Make(res, tmpfile.Name(), "")
	if err != nil {
		t.Fatalf("expected job model document to be generated. error: %v", err)
	}

	// check the file's content
	e := `
- name: random trigger on dummy-host
  description: random trigger on dummy-host
  nodefilters:
    filter: dummy-host
  sequence:
    commands:
    - exec: ""
- name: prefix random trigger on dummy-host2
  description: prefix random trigger on dummy-host2
  nodefilters:
    filter: dummy-host2
  sequence:
    commands:
    - exec: ""`

	c, err := ioutil.ReadFile(tmpfile.Name())

	if err != nil {
		t.Fatalf("expected file content to be %v\n error: %v", e, err)
	}

	content := strings.Trim(string(c), "\n")
	expectedContent := strings.Trim(e, "\n")

	if content != expectedContent {
		t.Errorf("test returned \n%v\n\n expected file content to be \n%v", content, expectedContent)
	}

	// Truncate temp file (We want previous file contents gone before next test)
	err = tmpfile.Truncate(0)
	if err != nil {
		t.Fatalf("expected file to be truncated. error: %v", err)
	}

	// With prefix for triggers
	err = Make(res, tmpfile.Name(), "prefix")
	if err != nil {
		t.Fatalf("expected job model document to be generated. error: %v", err)
	}

	// check the file's content
	e = `
- name: prefix random trigger on dummy-host2
  description: prefix random trigger on dummy-host2
  nodefilters:
    filter: dummy-host2
  sequence:
    commands:
    - exec: ""`

	c, err = ioutil.ReadFile(tmpfile.Name())

	if err != nil {
		t.Fatalf("expected file content to be %v\n error: %v", e, err)
	}

	content = strings.Trim(string(c), "\n")
	expectedContent = strings.Trim(e, "\n")

	if content != expectedContent {
		t.Errorf("test returned \n%v\n\n expected file content to be \n%v", content, expectedContent)
	}

}
