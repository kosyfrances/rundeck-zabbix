package jobs

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/kosyfrances/rundeck-zabbix/lib/utils"
	"github.com/kosyfrances/rundeck-zabbix/lib/zabbix"
)

type nodeFilters struct {
	Filter string `yaml:"filter"`
}

type commands struct {
	Exec string `yaml:"exec"`
}

type sequence struct {
	Commands []commands `yaml:"commands"`
}

type job struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	NodeFilters nodeFilters `yaml:"nodefilters"`
	Sequence    sequence    `yaml:"sequence"`
}

/*
Make generates Rundeck job model document in Yaml format
with triggers' information from Zabbix

The file output similar to:

- name: test-job
  description: Just a random test
  nodefilters:
    filter: 'dummy-host'
  sequence:
    commands:
    - exec: ""
*/
func Make(results zabbix.TriggerResults, filePath, prefix string) error {

	var jobList []job

	for _, result := range results {
		if strings.HasPrefix(result.Description, prefix) {
			for _, host := range result.Hosts {
				resultFilter := nodeFilters{
					Filter: host.Name,
				}

				seq := sequence{
					Commands: []commands{{Exec: ""}},
				}

				// Note: Rundeck job names does not accept slashes
				// https://rundeck.lighthouseapp.com/projects/59277/tickets/558
				// The user is responsible for ensuring that the trigger description on Zabbix does not contain slashes,
				// else even if it gets added to the job file, the user will not be able to load it in Rundeck.

				j := job{
					Name:        result.Description,
					Description: result.Description,
					NodeFilters: resultFilter,
					Sequence:    seq,
				}
				jobList = append(jobList, j)
			}
		}
	}

	d, err := yaml.Marshal(&jobList)
	if err != nil {
		return fmt.Errorf("cannot marshal resource to yaml. %v", err)
	}
	return utils.DumpToFile(filePath, d)
}
