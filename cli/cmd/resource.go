package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	URL := viper.GetString("zabbix.URL")
	key := viper.GetString("zabbix.api_key")
	a, err := zabbix.NewAPI(URL, key, "", "")
	if err != nil {
		log.Error("cannot find needed params.", err)
	}

	res, err := a.GetHostsInfo()
	if err != nil {
		log.Error("cannot get hosts info.\n", err)
		return
	}

	if len(res) == 0 {
		log.Warn("No resource found.")
		return
	}

	if err = resources.Make(res, filePath); err != nil {
		log.Error("cannot generate resource.\n", err)
	} else {
		log.Info("Resources generated in ", filePath)
	}
}
