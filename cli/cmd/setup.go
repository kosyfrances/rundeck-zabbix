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

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "initialise tool configuration",
	Long:  `initialise tool configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		var zabbixUrl, zabbixUser, zabbixPassword, rundeckUrl, rundeckApiKey string

		fmt.Println("Enter the Zabbix Server URL:")
		fmt.Scanln(&zabbixUrl)

		fmt.Println("Enter the Zabbix Username:")
		fmt.Scanln(&zabbixUser)

		fmt.Println("Enter the Zabbix Password:")
		fmt.Scanln(&zabbixPassword)

		fmt.Println("Enter the Rundeck URL:")
		fmt.Scanln(&rundeckUrl)

		fmt.Println("Enter the Rundeck API_KEY:")
		fmt.Scanln(&rundeckApiKey)

		zabbixConfig := lib.ZabbixConfig{
			Url:      zabbixUrl,
			UserName: zabbixUser,
			Password: zabbixPassword,
		}

		rundeckConfig := lib.RundeckConfig{
			Url:    rundeckUrl,
			ApiKey: rundeckApiKey,
		}

		newConfig := lib.Config{
			Zabbix:  zabbixConfig,
			Rundeck: rundeckConfig,
		}

		newConfig.Save()
		log.Info("Configuration file generated in " + lib.ConfigDirectory)
	},
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
