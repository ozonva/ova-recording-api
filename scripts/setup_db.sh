#!/bin/bash -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

source ${SCRIPT_DIR}/db_conf.sh

EXEC="sudo -u postgres psql -p ${DB_PORT}"

sudo pg_createcluster ${PG_VERSION} ${CLUSTER_NAME} -p ${DB_PORT}
sudo pg_ctlcluster ${PG_VERSION} ${CLUSTER_NAME} start

$EXEC "create database ${DB_NAME}"
$EXEC "create role ${DB_USER} with encrypted password '${DB_PASSWORD}'"
$EXEC "grant all on database ${DB_NAME} to ${DB_USER}"
$EXEC "alter role ${DB_USER} with login"
