package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/resources"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

var filePath string

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "generate Rundeck resources",
	Long: `generate Rundeck resource model document in Yaml format with hosts' information from Zabbix

	If a file is given, the generate resources are appended to the file.
	If the given filePath does not exist, it gets created.
	Otherwise, a resources.yml file is generated in the current path.`,
	Run: runResource,
}

func init() {
	generateCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVar(&filePath, "path", "resources.yml", "Path to file where generated Resources will be written")
}

func runResource(cmd *cobra.Command, args []string) {
	newConfig, err := lib.NewConfigFromFile(lib.ConfigPath)
	if err != nil {
		log.Errorf("cannot create config from file. %v", err)
		return
	}
	URL := newConfig.Zabbix.URL
	key := newConfig.Zabbix.APIKey

	a, err := zabbix.CreateClientUsingAPIKey(URL, key)
	if err != nil {
		log.Errorf("cannot find needed params. %v", err)
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
