#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

PG_VERSION=9.4
CLUSTER_NAME=ozonva

DB_PORT=5433
DB_HOST=localhost
DB_NAME=appointments
DB_USER=recording_user
DB_PASSWORD=recording_password
CONNECTION="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} database=${DB_NAME}"

MIGRATIONS_DIR="${SCRIPT_DIR}/../migrations"
