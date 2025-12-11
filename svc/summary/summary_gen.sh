#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)
CONFIG_FILE="${CUR}/config/tcm_summary.yaml"
OUT_FILE="${CUR}/summary.gen.go"
PKG_NAME="summary"

cd ${CUR}

echo "Generating summary/summary.gen.go"
go run ./summary_gen/gen.go -c ${CONFIG_FILE} -o ${OUT_FILE} -p ${PKG_NAME}
goimports -w ${OUT_FILE}
