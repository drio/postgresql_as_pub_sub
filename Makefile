all:
	psql $(DB_NAME)

var/postgres:
	initdb ./var/postgres

start-fg: var/postgres
	/opt/homebrew/opt/postgresql/bin/postgres -D ./var/postgres

start-fg-default:
	/opt/homebrew/opt/postgresql/bin/postgres es -D /opt/homebrew/var/postgres

DB_NAME=pub_sub_test

up-sql:
	echo "create database $(DB_NAME)" | psql
	cat sql/create.sql | psql pub_sub_test
	cat sql/trigger.sql | psql pub_sub_test

clean:
	psql < sql/clean.sql
	echo "drop database IF EXISTS $(DB_NAME)" | psql

new_job:
	cat sql/new_job.sql | psql $(DB_NAME)

list_jobs:
	echo "select * from ps_jobs" | psql $(DB_NAME)