package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kosyfrances/rundeck-zabbix/lib"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rundeck-zabbix",
	Short: "rundeck-zabbix cli tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to rundeck-zabbix config file")
}

// initConfig reads in config file.
func initConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		log.Error("no config file given.")
		os.Exit(1)
	}

	if !lib.FileExists(cfgFile) && setupCmd.CalledAs() == "" {
		log.Error("config file does not exist. run setup to create new configuration.")
		os.Exit(1)
	} else {
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			log.Info("Loading config file:", viper.ConfigFileUsed())
		}
	}
}
