package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ZabbixConfig struct holds Zabbix configuration info
type ZabbixConfig struct {
	URL      string `json:"url"`
	UserName string `json:"username"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
}

// RundeckConfig struct holds Zabbix configuration info
type RundeckConfig struct {
	URL    string `json:"url"`
	APIKey string `json:"api_key"`
}

// Config struct holds the setup config file info
type Config struct {
	Zabbix  ZabbixConfig  `json:"zabbix"`
	Rundeck RundeckConfig `json:"rundeck"`
}

// Save object config
func (object *Config) Save(configPath string) error {
	data, err := json.MarshalIndent(*object, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config. Error: %v", err)
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file. Error: %v", err)
	}
	return nil
}

// NewConfigFromFile creates new config from file
func NewConfigFromFile(configPath string) (*Config, error) {
	config := Config{}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file. Error: %v", err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal json. Error: %v", err)
	}
	return &config, nil
}

// FileExists checks if file exists and returns a bool value
func FileExists(name string) bool {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return true
	}
	return false
}
