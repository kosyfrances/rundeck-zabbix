package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib/resources"
)

// resourceCmd represents the resource command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "generate Rundeck resources",
	Long: `generate Rundeck resource model document in Yaml format with hosts' information from Zabbix

	If a file path is given, the generated resources are appended to the file.
	If the given file Path does not exist, it gets created.
	Otherwise, a resources.yml file is generated in the current path.`,
	Run: runResources,
}

func init() {
	generateCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().StringVar(&filePath, "file", "resources.yml", "Path to file where generated Resources will be written")
}

func runResources(cmd *cobra.Command, args []string) {
	a, err := createZabbixClient()
	if err != nil {
		log.Errorf("error creating Zabbix client. %v", err)
		return
	}

	res, err := a.GetHostsInfo()
	if err != nil {
		log.Errorf("cannot get hosts info. %v", err)
		return
	}

	if len(res) == 0 {
		// No servers are on Zabbix
		log.Warn("No resource found.")
		return
	}

	if err = resources.Make(res, filePath); err != nil {
		log.Errorf("cannot generate resource. %v", err)
	} else {
		log.Infof("Resources generated in %v", filePath)
	}
}
