package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

// Use for generating resources
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Use for generating utilities",
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func createZabbixClient() (*zabbix.API, error) {
	newConfig, err := lib.NewConfigFromFile(lib.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("cannot create config from file. %v", err)

	}
	URL := newConfig.Zabbix.URL
	key := newConfig.Zabbix.APIKey

	return zabbix.CreateClientUsingAPIKey(URL, key)
}
