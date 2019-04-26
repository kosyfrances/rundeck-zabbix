package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kosyfrances/rundeck-zabbix/lib/utils"
)

// BuildRundeckURLEndpoint builds a URL endpoint given a baseURL and an endpoint.
// It returns the full URL endpoint and an error if any.
func BuildRundeckURLEndpoint(baseURL, endpoint string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("cannot parse Rundeck base URL")
	}

	queryEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("cannot parse Rundeck query Endpoint")
	}

	return base.ResolveReference(queryEndpoint).String(), nil
}

/*
GetRundeckJobID gets the job ID from Rundeck API.

Note:
Rundeck API does not return Host/Node name alongside when you try to get a job.
Here is a scenario:
A trigger event happens on Zabbix, we collect the Trigger name and Host name,
make an API call to Rundeck and try to get the corresponding job name on Rundeck.
Rundeck returns list of jobs with that exact name
(assuming there are multiple jobs with same name for different servers).

A sample response is:
[{"id":"91c755ac-dab2-475c-b903-92e5ff59e740","name":"[RD] etc passwd has been changed on {HOST.NAME}","group":null,"project":"zabideck","description":"[RD] /etc/passwd has been changed on {HOST.NAME}","href":"http://localhost:4440/api/27/job/91c755ac-dab2-475c-b903-92e5ff59e740","permalink":"http://localhost:4440/project/zabideck/job/show/91c755ac-dab2-475c-b903-92e5ff59e740","scheduled":false,"scheduleEnabled":true,"enabled":true},
{"id":"826e88dc-926b-415f-a724-a3a847aae3f2","name":"[RD] etc passwd has been changed on {HOST.NAME}","group":null,"project":"zabideck","description":"[RD] /etc/passwd has bee
n changed on {HOST.NAME}","href":"http://localhost:4440/api/27/job/826e88dc-926b-415f-a724-a3a847aae3f2","permalink":"http://localhost:4440/project/zabideck/job/show/826e88dc-926b-415f-a724-a3a847aae3f2","scheduled":false,"scheduleEnabled":true,"enabled":true}]

How do we figure out what exact job to run in this case?

This function is implemented with the assumption that the job name (i.e Trigger name)
being given is unique per project, else it will return the first match on the list.
*/
func GetRundeckJobID(URL string) (string, error) {
	type response struct {
		ID string `json:"id"`
	}

	var r []response

	resp, err := utils.MakeRundeckRequest(http.MethodGet, URL, nil)
	if err != nil {
		return "", fmt.Errorf("cannot make Rundeck API call to get job ID. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", fmt.Errorf("cannot decode response. Error: %v", err)
	}

	if len(r) == 0 {
		return "", fmt.Errorf("no matching jobs found")
	}

	return r[0].ID, nil
}

type job struct {
	Name string `json:"name"`
}

// JobExecResponse struct to hold result from Rundeck job execution
type JobExecResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Project string `json:"project"`
	Job     job    `json:"job"`
}

// ExecuteRundeckJob executes a job on Rundeck given a URL endpoint.
// It returns an error if any.
func ExecuteRundeckJob(URL string) (*JobExecResponse, error) {

	r := &JobExecResponse{}

	resp, err := utils.MakeRundeckRequest(http.MethodPost, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot make Rundeck API call to execute job. Error: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response. Error: %v", err)
	}

	return r, nil
}
