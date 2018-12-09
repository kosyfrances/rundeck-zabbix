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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var z_url, z_user, z_password, r_url, r_api_key string

		fmt.Println("Enter the Zabbix Server URL:")
		fmt.Scanln(&z_url)

		fmt.Println("Enter the Zabbix Username:")
		fmt.Scanln(&z_user)

		fmt.Println("Enter the Zabbix Password:")
		fmt.Scanln(&z_password)

		fmt.Println("Enter the Rundeck URL:")
		fmt.Scanln(&r_url)

		fmt.Println("Enter the Rundeck API_KEY:")
		fmt.Scanln(&r_api_key)

		zabbixConfig := lib.ZabbixConfig{
			Url:      z_url,
			UserName: z_user,
			Password: z_password,
		}

		rundeckConfig := lib.RundeckConfig{
			Url:    r_url,
			ApiKey: r_api_key,
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
