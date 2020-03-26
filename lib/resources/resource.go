package resources

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/kosyfrances/rundeck-zabbix/lib"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

// Resource struct holds Rundeck resource document mapping info
type resource struct {
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
func Make(results zabbix.HostResults, file string) error {
	m := make(map[string]resource)

	for _, result := range results {
		m[result.Name] = resource{
			Name:        result.Host,
			Node:        result.Name,
			Description: result.Description,
		}
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		return fmt.Errorf("cannot marshal resource to yaml. %v", err)
	}
	return lib.DumpToFile(file, d)
}
