package cmd

import (
	"fmt"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

func createZabbixClient() (*zabbix.API, error) {
	newConfig, err := lib.NewConfigFromFile(lib.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("cannot create config from file. %v", err)

	}
	URL := newConfig.Zabbix.URL
	key := newConfig.Zabbix.APIKey

	return zabbix.CreateClientUsingAPIKey(URL, key)
}
