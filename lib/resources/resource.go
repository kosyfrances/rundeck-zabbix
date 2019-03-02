package resources

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

// Resource struct holds Rundeck resource document mapping info
type Resource struct {
	Name        string `yaml:"hostname"`
	Node        string `yaml:"nodename"`
	Description string `yaml:"description"`
	OSArch      string `yaml:"osArch"`
	OSFamily    string `yaml:"osFamily"`
	OSName      string `yaml:"osName"`
	OSVersion   string `yaml:"osVersion"`
	Tags        string `yaml:"tags"`
	Username    string `yaml:"username"`
	SSHKeypath  string `yaml:"ssh-keypath"`
}

/*
Make generates Rundeck resource model document in Yaml format
with hosts' information from Zabbix

The file output similar to:

dummy-host:
  hostname: dummy-host
  nodename: dummy-host
  description: ""
  osArch: ""
  osFamily: ""
  osName: ""
  osVersion: ""
  tags: ""
  username: ""
  ssh-keypath: ""
*/
func Make(results zabbix.Results, file string) error {
	m := make(map[string]Resource)

	for _, result := range results {
		m[result.Name] = Resource{
			Name:        result.Host,
			Node:        result.Name,
			Description: result.Description,
		}
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return fmt.Errorf("cannot marshal resource to yaml. %v", err)
	}
	return dumpToFile(file, d)
}

func dumpToFile(filePath string, data []byte) error {
	// If the filePath doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file. %v", err)
	}

	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("cannot write to file. %v", err)
	}

	return nil
}
