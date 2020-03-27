build-zabbix-server:
	./dev/zabbix-server.sh

build-rundeck:
	./dev/rundeck.sh

build-dummy-host:
	./dev/dummy-host.sh
	docker start dummy-host

build-containers: build-zabbix-server build-dummy-host build-rundeck

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

start-containers: start-zabbix-server start-rundeck start-dummy-host

stop-containers: stop-zabbix-server stop-rundeck stop-dummy-host

test:
	golint `go list ./... | grep -v /vendor/`
	@echo ""
	go test ./...

build-linux:
	pushd cli/ && \
	CGO_ENABLED=0 GOOS=linux go build -o ../rundeck-zabbix  && \
	popd

build:
	pushd cli/ && \
	go build -o ../rundeck-zabbix  && \
	popd
