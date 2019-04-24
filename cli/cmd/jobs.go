package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib/jobs"
)

var prefix string
var jobsFilePath string

// jobsCmd represents the jobs command for generating Rundeck jobs
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "generate Rundeck jobs",
	Long: `generate Rundeck jobs document in Yaml format with triggers information about hosts from Zabbix

	If a prefix is given, jobs will be generated from triggers with the given prefix.
	Otherwise, jobs will be generated from all the triggers in Zabbix.

	If a file is given, the generated jobs are appended to the file.
	If the given file path does not exist, it gets created.
	Otherwise, a jobs.yml file is generated in the current path.
	Note that Rundeck job names must not contain slashes.
	This means that the Trigger name from Zabbix must not also contain slashes.`,
	Run: generateJobs,
}

func init() {
	generateCmd.AddCommand(jobsCmd)
	jobsCmd.Flags().StringVar(&jobsFilePath, "file", "jobs.yml", "Path to file where generated Jobs will be written.")
	jobsCmd.Flags().StringVar(&prefix, "prefix", "", "Generate triggers from the given prefix, otherwise, jobs will be generated from all the triggers in Zabbix.")
}

func generateJobs(cmd *cobra.Command, args []string) {
	a, err := createZabbixClient()
	if err != nil {
		log.Errorf("error creating Zabbix client. %v", err)
		return
	}

	res, err := a.GetTriggersInfo()
	if err != nil {
		log.Errorf("cannot get triggers info. %v", err)
		return
	}

	// No active triggers are on Zabbix
	if len(res) == 0 {
		log.Warn("No active triggers found.")
		return
	}

	if err = jobs.Make(res, jobsFilePath, prefix); err != nil {
		log.Errorf("cannot generate jobs. %v", err)
	} else {
		log.Infof("Jobs generated in %v", jobsFilePath)
	}
}
