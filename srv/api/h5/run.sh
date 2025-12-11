#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=`dirname $0`

cd ${CUR}

SERVER_ADDR=0.0.0.0:9092
ETCD_ADDR=localhost:2379

# 启动服务
go run . \
    --log_level=DEBUG \
    --server_address=${SERVER_ADDR} \
    --config_etcd_address=${ETCD_ADDR} \
    # --enable_jwt

# go run . \
#     --registry=etcd \
#     --log_level=DEBUG \
#     --register_interval=5 \
#     --register_ttl=10 \
#     --server_address=${SERVER_ADDR} \
#     --client_pool_size=100 \
#     --config_etcd_address=${ETCD_ADDR} \
#     # --enable_jwt
