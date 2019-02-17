package cmd

import (
	"fmt"

	"github.com/prometheus/common/log"
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
	URL := viper.GetString("zabbix.url")
	key := viper.GetString("zabbix.api_key")
	a := zabbix.NewAPI(URL, key)

	res, err := a.GetHostsInfo()
	if err != nil {
		log.Error("cannot get hosts info.\n", err)
	}

	if err = resources.Make(res, filePath); err != nil {
		fmt.Println(err)
		log.Error("cannot generate resource.\n", err)
	}
}
