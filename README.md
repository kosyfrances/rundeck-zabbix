[![Build Status](https://travis-ci.org/kosyfrances/rundeck-zabbix.svg?branch=master)](https://travis-ci.org/kosyfrances/rundeck-zabbix)

# rundeck-zabbix
rundeck-zabbix is a Runbook automation tool, that integrates Zabbix trigger-based actions into Rundeck. It is based on https://kosyfrances.github.io/rundeck-zabbix/ and built with Zabbix API v4.2 and Rundeck API v30.

**Basic flow:**

* A service stops running
* Zabbix fires trigger
* Zabbix action calls rundeck-zabbix
* rundeck-zabbix executes job on Rundeck
* rundeck-zabbix sends acknowledgement to Zabbix
* Zabbix receives acknowledgement
* Ops continues partying, nothing left for them to do

### Prerequisites
* [Zabbix server](https://www.zabbix.com/download)
* [Rundeck server](https://www.rundeck.com/open-source/download)
* [Go 1.13.x](https://golang.org/doc/install)

### Installation
Get the [latest release](https://github.com/kosyfrances/rundeck-zabbix/releases), unarchive and build with:
```bash
# For Linux OS
make build

# For generic build specific to your machine's OS
make build-generic
```

There is more information [here](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies) on how to compile Go packages and dependencies.

Run `rundeck-zabbix help` to get list of available commands.

### Initialise tool configuration
You'd need the following information:
* Zabbix Server URL in the form of `https://ZABBIX_URL/api_jsonrpc.php`
* Zabbix username and password (required to generate Zabbix API key)
* Rundeck URL
* Rundeck API key (usually gotten from `https://RUNDECK_URL:4440/user/profile`)

```bash
rundeck-zabbix setup --config=/path/to/dir/.config.json
```
where `/path/to/dir` already exists.

Logs can be found in `/tmp/rundeck-zabbix.log`.

### Generate Rundeck resources template
Rundeck uses a resource document to declare the resource models used by a project to define the set of Nodes that are available. This file is usually found in `/etc/rundeck/projects/project_name/` or `/var/rundeck/projects/project_name` in the Rundeck server.

To map Zabbix hosts to Rundeck nodes, run the following to generate a `resources.yml` template and modify as appropriate. You can read more about [node entries here](https://docs.rundeck.com/docs/man5/resource-yaml.html) and format yours appropriately.

```bash
rundeck-zabbix generate resources --config=/path/to/dir/.config.json
```
Do not forget to add the existing Rundeck’s `localhost` resource definition gotten from the original resource document file in the server. The localhost resource definition looks like this:

```
localhost:
  description: Rundeck server
  hostname: localhost
  nodename: localhost
  osArch: amd64
  osFamily: unix
  osName: Linux
  osVersion: 4.4.0-66-generic
  tags: 'localhost'
  username: rundeck
```

Copy the file over to `/etc/rundeck/projects/project_name/resources.yml` or `/var/rundeck/projects/project_name/resources.yml` on Rundeck’s server.

### Generate Rundeck jobs template
To map Zabbix triggers to Rundeck jobs, run the following to generate a `jobs.yml` template and modify as appropriate. Update the `commands` section with the appropriate command you want Rundeck to execute. You can read more about [job entries here](https://docs.rundeck.com/docs/man5/job-yaml.html) and format yours appropriately.

```bash
rundeck-zabbix generate jobs --config=/path/to/dir/.config.json
```
Next, load the [jobs file to Rundeck](https://docs.rundeck.com/docs/man5/job-yaml.html#loading-and-unloading).

### Run middleware
The remote command for [Zabbix action](https://www.zabbix.com/documentation/4.2/manual/config/notifications/action/operation/remote_command) will be
```
rundeck-zabbix run job --project=RUNDECK-PROJECT --trigger={TRIGGER.NAME} --event={EVENT.ID} --config="/path/to/.config.json"
```

![Photo of Zabbix Action Page](/dev/assets/zabbix-action-page.png)

### Developing locally
Refer to the [dev setup guide](/dev/README.md).

### Contributing
Refer to the [contributing guide](/CONTRIBUTING.md).

### Versioning
We use [SemVer](https://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kosyfrances/rundeck-zabbix/tags).

### Authors
See the list of [contributors](https://github.com/kosyfrances/rundeck-zabbix/graphs/contributors) who participated in this project.

### License
This project is licensed under the MIT License - see the [LICENSE](/LICENSE) file for details

### Acknowledgments
Special thanks to [Femi](https://github.com/osule) and [Oscar](https://github.com/0sc) for helping out with code reviews.
