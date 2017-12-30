#!/bin/bash

# This script spins up rundeck in docker
# https://hub.docker.com/r/jordan/rundeck/


### Note: You can access Rundeck at http://localhost:4440/user/login
### Username:admin, password: admin

set -eo pipefail

docker run -p 4440:4440 \
    -e EXTERNAL_SERVER_URL=http://localhost:4440 \
    -e RUNDECK_PASSWORD="rundeck" \
    --name rundeck \
    -t jordan/rundeck:latest
