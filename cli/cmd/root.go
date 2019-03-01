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
	Use:   "cli",
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// setup config file with defaults
	if lib.FileExists(lib.ConfigPath) == false {
		// only warn when not running setup command
		log.Warn("Configuration file not found. Please setup using `setup` command")
	} else {
		viper.SetConfigName(lib.AppConfigName)   // name of config file (without extension)
		viper.AddConfigPath(lib.ConfigDirectory) // adding home directory as first search path
		viper.AutomaticEnv()                     // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			log.Info("Loading config file:", viper.ConfigFileUsed())
		}

	}

}
