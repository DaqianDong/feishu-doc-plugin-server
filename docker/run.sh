#!/bin/bash
# Usage="Usage: $0 version"

set -x
PA_DOCKER_REGISTRY_HOST='registry.sofunny.io/project-tr'
version='1.0'
PA_SERVER_NAME='http://10.30.40.176'
PA_PORT=9004
PA_WEB_PORT=3001
PA_DB_HOST='10.30.40.176'
PA_DB_USER='root'
PA_DB_PWD='1q2w3e4r'
PA_DB_NAME='tr_pa'
PA_DB_PORT=3306
PA_ETCD='10.30.40.176:20034'
PA_GAMEDB_PARAMS='root:1q2w3e4r@tcp(10.30.40.176:3306)/tr_game'
PA_PPSERVER_DB_PARAMS='root:1q2w3e4r@tcp(10.30.40.176:3306)/tr_passport'
PA_APP_ID='123456'
PA_GM_API_ADDR='http://127.0.0.1:7202'
TA_URL='http://localhost:8992?token=123456&projectId=1'

docker rm -f tr_pa
docker run -d -p 3001:3001 \
  -e SERVER_NAME=$PA_SERVER_NAME \
  -e PORT=$PA_PORT \
  -e WEB_PORT=$PA_WEB_PORT \
  -e DB_HOST=$PA_DB_HOST \
  -e DB_USER=$PA_DB_USER \
  -e DB_PASS=$PA_DB_PWD \
  -e DB_NAME=$PA_DB_NAME \
  -e DB_PORT=$PA_DB_PORT \
  -e ETCD=$PA_ETCD \
  -e GAMEDB_PARAMS=$PA_GAMEDB_PARAMS \
  -e PPSERVER_DB_PARAMS=$PA_PPSERVER_DB_PARAMS \
  -e APP_ID=$PA_APP_ID \
  -e GM_API_ADDR=$PA_GM_API_ADDR \
  -e TA_URL=$TA_URL \
  -e VERSION=$version \
  -v ~/Data/logs/tr_pa:/data/platform_admin/data/log  \
  --name tr_pa \
  $PA_DOCKER_REGISTRY_HOST/platform_admin:$version
