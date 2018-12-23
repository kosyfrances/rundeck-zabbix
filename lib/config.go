package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ZabbixConfig struct {
	Url      string `json:"url"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type RundeckConfig struct {
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
}

type Config struct {
	Zabbix  ZabbixConfig  `json:"zabbix"`
	Rundeck RundeckConfig `json:"rundeck"`
}

var (
	UserHome        string
	ConfigPath      string
	ConfigDirectory string
)

const APP_CONFIG_DIRECTORY = "/.rundeck-zabbix/"
const APP_CONFIG_NAME = "config"
const APP_CONFIG_EXTENSION = "json"

func init() {
	UserHome = os.ExpandEnv("$HOME")
	ConfigDirectory = UserHome + APP_CONFIG_DIRECTORY
	ConfigPath = ConfigDirectory + APP_CONFIG_NAME + "." + APP_CONFIG_EXTENSION
	if FileExists(ConfigDirectory) == false {
		os.MkdirAll(ConfigDirectory, os.ModePerm)
	}
}

func (object *Config) Save() {
	data, err := json.MarshalIndent(*object, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(ConfigPath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func NewConfigFromFile() *Config {
	config := Config{}
	b, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
