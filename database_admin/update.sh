#!/bin/bash

export PGHOST=$DB_HOST
export PGUSER=$DB_ADMIN_USER
export PGPASSWORD=$DB_ADMIN_PASSWD
export PGDATABASE=$DB_NAME
export PGPORT=$DB_PORT
export PGSSLMODE=$DB_SSLMODE
export PGSSLROOTCERT=$DB_SSLROOTCERT

DB_USER=$DB_ADMIN_USER DB_PASSWD=$DB_ADMIN_PASSWD WAIT_FOR_EMPTY_DB=1 ./scripts/wait-for-services.sh

if [[ $RESET_SCHEMA == "true" ]]; then
  psql -c "DROP SCHEMA IF EXISTS public CASCADE"
  psql -c "CREATE SCHEMA IF NOT EXISTS public"
  psql -c "GRANT ALL ON SCHEMA public TO ${DB_USER}"
  psql -c "GRANT ALL ON SCHEMA public TO public"
fi

set -e -o pipefail

# we cain either create the database from scratch, or upgrade running database

# Create users if they don't exist
echo "Creating application components users"
psql -f ./database_admin/schema/create_users.sql
psql -f ./database_admin/schema/create_schema.sql

#Fail on unset passwords
set -u

echo "Setting user passwords"
# Set specific password for each user. If the users are already created, change the password.
# This is performed on each startup in order to ensure users have latest pasword
psql -c "ALTER USER vmaas_writer WITH PASSWORD '${WRITER_PASSWORD}'"
psql -c "ALTER USER vmaas_reader WITH PASSWORD '${READER_PASSWORD}'"
