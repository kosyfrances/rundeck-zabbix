build-zabbix-server:
	./dev/zabbix-server.sh

build-rundeck:
	./dev/rundeck.sh

build-dummy-host:
	./dev/dummy-host.sh
	docker start dummy-host

build-all: build-zabbix-server build-dummy-host build-rundeck

start-zabbix-server:
	echo "Starting Zabbix server containers"
	docker start mysql-server
	docker start zabbix-java-gateway
	docker start zabbix-server-mysql
	docker start zabbix-web-nginx-mysql

stop-zabbix-server:
	echo "Stopping Zabbix server containers"
	docker stop zabbix-web-nginx-mysql
	docker stop zabbix-server-mysql
	docker stop zabbix-java-gateway
	docker stop mysql-server

start-rundeck:
	echo "Starting rundeck container"
	docker start rundeck

stop-rundeck:
	echo "Stopping rundeck container"
	docker stop rundeck

start-dummy-host:
	echo "Starting Zabbix agent"
	docker start dummy-host

stop-dummy-host:
	echo "Stopping Dummy host"
	docker stop dummy-host

start-all: start-zabbix-server start-rundeck start-dummy-host

stop-all: stop-zabbix-server stop-rundeck stop-dummy-host
