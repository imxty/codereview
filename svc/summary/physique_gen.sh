#!/usr/bin/env bash
set -e
set -u
set -o pipefail

CUR=$(dirname $0)
CONFIG_FILE="${CUR}/config/physique.yaml"
OUT_FILE="${CUR}/physique.gen.go"
PKG_NAME="summary"

cd ${CUR}

echo "Generating physique/physique.gen.go"
go run ./physique_gen/gen.go -c ${CONFIG_FILE} -o ${OUT_FILE} -p ${PKG_NAME}
goimports -w ${OUT_FILE}
