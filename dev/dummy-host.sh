#!/bin/bash

# This script builds and runs a docker container with apache to serve as a Zabbix host

docker build -t dummy-host:latest -f $(PWD)/dev/Dockerfile .
docker run --name dummy-host -t -p 8001:80 -p 20050:10050 -d dummy-host:latest
