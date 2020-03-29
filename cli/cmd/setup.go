package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "initialise tool configuration",
	Run:   runSetup,
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
}

func runSetup(cmd *cobra.Command, args []string) {
	cfg, err := lib.NewConfigFromFile(cfgFile)
	if err != nil {
		log.Infof("creating configuration ...")
		cfg = &lib.Config{}
	}

	// Read Zabbix config

	fmt.Println(appendValueToPrompt("Enter the Zabbix Server URL:", cfg.Zabbix.URL))
	fmt.Scanln(&cfg.Zabbix.URL)

	fmt.Println(appendValueToPrompt("Enter the Zabbix Username:", cfg.Zabbix.UserName))
	fmt.Scanln(&cfg.Zabbix.UserName)

	fmt.Println(appendValueToPrompt("Enter the Zabbix Password:", cfg.Zabbix.Password))
	fmt.Scanln(&cfg.Zabbix.Password)

	// Read Rundeck config

	fmt.Println(appendValueToPrompt("Enter the Rundeck URL:", cfg.Rundeck.URL))
	fmt.Scanln(&cfg.Rundeck.URL)

	fmt.Println(appendValueToPrompt("Enter the Rundeck API_KEY:", cfg.Rundeck.APIKey))
	fmt.Scanln(&cfg.Rundeck.APIKey)

	// Generate Zabbix API key
	z, err := zabbix.CreateClientUsingAuth(cfg.Zabbix.URL, cfg.Zabbix.UserName, cfg.Zabbix.Password)
	if err != nil {
		log.Errorf("cannot create Zabbix API client. %v", err)
		return
	}

	cfg.Zabbix.APIKey, err = z.GetKey(timeout)
	if err != nil {
		log.Errorf("cannot generate Zabbix API key. %v", err)
		return
	}

	err = cfg.Save(cfgFile)
	if err != nil {
		log.Errorf("cannot save configuration. %v", err)
	} else {
		log.Infof("Configuration file generated in %v", cfgFile)
	}
}
