# Dev setup

### Zabbix setup on docker
* Run `./zabbix.sh` to spin up containers needed for Zabbix to run on docker.
* Visit `http://localhost:80` to access Zabbix web interface. Username: `Admin`, password: `zabbix`

### Rundeck setup on docker
* Run `./rundeck.sh` to spin up Rundeck container.
* Visit `http://localhost:4440/user/login` to access Rundeck web interface. Username: `admin`, password: `admin`

### Container management
* Stop Zabbix containers --> `./container.sh stop_zabbix`
* Start Zabbix containers -->  `./container.sh start_zabbix`
* Stop Rundeck containers --> `./container.sh stop_rundeck`
* Start Rundeck containers --> `./container.sh start_rundeck`
