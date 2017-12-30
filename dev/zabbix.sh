#!/bin/bash

# This script runs docker containers containing Zabbix server with MySQL database support,
# Zabbix web interface based on the Nginx web server and Zabbix Java gateway.
# https://www.zabbix.com/documentation/3.4/manual/installation/containers


### NOTE: You can access Zabbix at http://localhost:80
### Zabbix web username: Admin, password: zabbix

set -eo pipefail

# Start empty MySQL server instance
docker run --name mysql-server -t \
      -e MYSQL_DATABASE="zabbix" \
      -e MYSQL_USER="zabbix" \
      -e MYSQL_PASSWORD="zabbix_pwd" \
      -e MYSQL_ROOT_PASSWORD="root_pwd" \
      -d mysql:5.7 \
      --character-set-server=utf8 --collation-server=utf8_bin

# Start Zabbix Java gateway instance
docker run --name zabbix-java-gateway -t \
      -d zabbix/zabbix-java-gateway:latest

# Start Zabbix server instance and link the instance with created MySQL server instance
docker run --name zabbix-server-mysql -t \
      -e DB_SERVER_HOST="mysql-server" \
      -e MYSQL_DATABASE="zabbix" \
      -e MYSQL_USER="zabbix" \
      -e MYSQL_PASSWORD="zabbix_pwd" \
      -e MYSQL_ROOT_PASSWORD="root_pwd" \
      -e ZBX_JAVAGATEWAY="zabbix-java-gateway" \
      --link mysql-server:mysql \
      --link zabbix-java-gateway:zabbix-java-gateway \
      -p 10051:10051 \
      -d zabbix/zabbix-server-mysql:latest

# Note: Zabbix server instance exposes 10051/TCP port (Zabbix trapper) to host machine.

# Start Zabbix web interface and link the instance with created MySQL server and Zabbix server instances
# Zabbix web interface instance exposes 80/TCP port (HTTP) to host machine.
docker run --name zabbix-web-nginx-mysql -t \
      -e DB_SERVER_HOST="mysql-server" \
      -e MYSQL_DATABASE="zabbix" \
      -e MYSQL_USER="zabbix" \
      -e MYSQL_PASSWORD="zabbix_pwd" \
      -e MYSQL_ROOT_PASSWORD="root_pwd" \
      --link mysql-server:mysql \
      --link zabbix-server-mysql:zabbix-server \
      -p 80:80 \
      -d zabbix/zabbix-web-nginx-mysql:latest
