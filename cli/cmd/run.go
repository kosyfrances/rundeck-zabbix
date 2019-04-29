package cmd

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/middleware"
)

var rundeckProject string
var zabbixTrigger string
var eventID string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run commands",
}

// jobCmd represents the job command for executing a Rundeck job from Zabbix server
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "run/execute a job in a Rundeck project",
	Long: `run/execute a job in a Rundeck project,
	given a Zabbix trigger name that matches the Rundeck job name to be executed.
	This is with the assumption that the job name (i.e Trigger name) being given is unique per project,
	else it will execute the first match on the list.
	Note that Rundeck job names must not contain slashes.
	This means that the Trigger name from Zabbix must not also contain slashes.
	`,
	Run: runJob,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.AddCommand(jobCmd)
	jobCmd.Flags().StringVar(&rundeckProject, "project", "", "Name of Rundeck project whose job will be executed.")
	jobCmd.Flags().StringVar(&zabbixTrigger, "trigger", "", "Name of Zabbix trigger that is an exact match of a Rundeck job name to be executed.")
	jobCmd.Flags().StringVar(&eventID, "event", "", "Zabbix Event ID to be acknowledged.")
}

func runJob(cmd *cobra.Command, args []string) {
	// Get Rundeck Auth Token
	newConfig, err := lib.NewConfigFromFile(cfgFile)
	if err != nil {
		log.Errorf("cannot create config from file. %v", err)
		return
	}
	authToken := newConfig.Rundeck.APIKey
	URL := newConfig.Rundeck.URL

	// Get job
	jobFilter := url.QueryEscape(zabbixTrigger)
	jobGetEndpoint := fmt.Sprintf("api/17/project/%s/jobs?authtoken=%s&jobExactFilter=%s", rundeckProject, authToken, jobFilter)
	URLEndpoint, err := middleware.BuildRundeckURLEndpoint(URL, jobGetEndpoint)
	if err != nil {
		log.Errorf("cannot build Rundeck URL job get endpoint. %v", err)
		return
	}

	ID, err := middleware.GetRundeckJobID(URLEndpoint)
	if err != nil {
		log.Errorf("cannot get Rundeck job ID. %v", err)
		return
	}

	// Run job
	jobRunEndpoint := fmt.Sprintf("/api/18/job/%s/run?authtoken=%s", ID, authToken)
	URLEndpoint, err = middleware.BuildRundeckURLEndpoint(URL, jobRunEndpoint)
	if err != nil {
		log.Errorf("cannot build Rundeck URL job run endpoint. %v", err)
		return
	}

	// Execute job
	execResp, err := middleware.ExecuteRundeckJob(URLEndpoint)
	if err != nil {
		log.Errorf("cannot execute Rundeck job. Job ID: %s; Error: %v", ID, err)
		return
	}
	log.Infof("Successfully executed Rundeck job. Job ID: %s", ID)

	// Send acknowledgment to Zabbix
	a, err := createZabbixClient()
	if err != nil {
		log.Errorf("error creating Zabbix client. %v", err)
		return
	}

	message := fmt.Sprintf(
		"Rundeck-Execution ID: %d, Status: %s, Job: %s, Project: %s", execResp.ID, execResp.Status, execResp.Job.Name, execResp.Project,
	)
	ackIDs, err := a.AcknowledgeEvent(eventID, message)
	if err != nil {
		log.Errorf("cannot acknowledge Zabbix Event. %v", err)
		return
	}

	if len(ackIDs) == 0 {
		log.Errorf("no Zabbix event acknowledged. %v", err)
		return
	}
	log.Infof("Acknowledged event on Zabbix. Event ID(s): %v", ackIDs)
}
