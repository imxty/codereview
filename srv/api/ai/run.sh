#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)

cd ${CUR}

SERVER_ADDR=0.0.0.0:9092
ETCD_ADDR=localhost:2379

# 启动服务
go run . \
    --server_address=${SERVER_ADDR} \
    --web_base_path="/" \
    --log_level=DEBUG
# --config_etcd_address=${ETCD_ADDR} \
