#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)
CONFIG_FILE="${CUR}/config/risk.yaml"
OUT_FILE="${CUR}/risk.gen.go"
PKG_NAME="summary"

cd ${CUR}

echo "Generating risk/risk.gen.go"
go run ./risk_gen/gen.go -c ${CONFIG_FILE} -o ${OUT_FILE} -p ${PKG_NAME}
goimports -w ${OUT_FILE}
