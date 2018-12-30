// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "initialise tool configuration",
	Long:  `initialise tool configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		var zabbixURL, zabbixUser, zabbixPassword, rundeckURL, rundeckAPIKey string
		zabbixConfigKeyName, rundeckConfigKeyName := lib.CONFIG_ZABBIX_KEY_NAME, lib.CONFIG_RUNDECK_KEY_NAME
		if viper.IsSet(zabbixConfigKeyName) {
			zabbixURL = getConfigValueByKey(zabbixConfigKeyName + "." + lib.CONFIG_URL_KEY_NAME)
			zabbixUser = getConfigValueByKey(zabbixConfigKeyName + "." + lib.CONFIG_USERNAME_KEY_NAME)
			zabbixPassword = getConfigValueByKey(zabbixConfigKeyName + "." + lib.CONFIG_PASSWORD_KEY_NAME)
		}

		fmt.Println(appendValueToPrompt("Enter the Zabbix Server URL:", zabbixURL))
		fmt.Scanln(&zabbixURL)

		fmt.Println(appendValueToPrompt("Enter the Zabbix Username:", zabbixUser))
		fmt.Scanln(&zabbixUser)

		fmt.Println(appendValueToPrompt("Enter the Zabbix Password:", zabbixPassword))
		fmt.Scanln(&zabbixPassword)

		if viper.IsSet(rundeckConfigKeyName) {
			rundeckURL = getConfigValueByKey(rundeckConfigKeyName + "." + lib.CONFIG_URL_KEY_NAME)
			rundeckAPIKey = getConfigValueByKey(rundeckConfigKeyName + "." + lib.CONFIG_API_KEY_KEY_NAME)
		}
		fmt.Println(appendValueToPrompt("Enter the Rundeck URL:", rundeckURL))
		fmt.Scanln(&rundeckURL)

		fmt.Println(appendValueToPrompt("Enter the Rundeck API_KEY:", rundeckAPIKey))
		fmt.Scanln(&rundeckAPIKey)

		zabbixConfig := lib.ZabbixConfig{
			Url:      zabbixURL,
			UserName: zabbixUser,
			Password: zabbixPassword,
		}

		rundeckConfig := lib.RundeckConfig{
			Url:    rundeckURL,
			ApiKey: rundeckAPIKey,
		}

		newConfig := lib.Config{
			Zabbix:  zabbixConfig,
			Rundeck: rundeckConfig,
		}

		newConfig.Save()
		log.Info("Configuration file generated in " + lib.ConfigDirectory)
	},
}

// gets value from config file loaded via viper
func getConfigValueByKey(key string) string {
	var result string
	value := viper.Get(key)
	if value != nil {
		result = value.(string)
	}
	return result
}

// mask values retrieved from config
//https://play.golang.org/p/E4k4zT9ATK
func mask(value string) string {
	if len(value) <= 4 {
		return value
	}
	return strings.Repeat("*", len(value)-4) + value[len(value)-4:]
}

func appendValueToPrompt(promptString, value string) string {
	if value != "" {
		return promptString + " [" + mask(value) + "]"
	}
	return promptString
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
