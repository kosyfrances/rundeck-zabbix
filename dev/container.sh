# Quick script to manage containers used in this project

# stop_rundeck --> Stop Rundeck container
# start_rundeck --> Start Rundeck container

# stop_zabbix --> Stop all Zabbix containers
# start_zabbix --> Start all Zabbix containers

set -eo pipefail

case "$1" in
"stop_rundeck")
    echo "Stopping rundeck container"
    docker stop rundeck
    ;;

"start_rundeck")
    echo "Starting rundeck container"
    docker start rundeck
    ;;

"stop_zabbix")
    echo "Stopping Zabbix containers"
    docker stop zabbix-web-nginx-mysql
    docker stop zabbix-server-mysql
    docker stop zabbix-java-gateway
    docker stop mysql-server
    ;;

"start_zabbix")
    echo "Starting Zabbix containers"
    docker start mysql-server
    docker start zabbix-java-gateway
    docker start zabbix-server-mysql
    docker start zabbix-web-nginx-mysql
    ;;

*)
    echo "I have no idea what you are trying to do. Please see README for available commands"
esac
