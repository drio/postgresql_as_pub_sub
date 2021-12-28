select:
	echo "select * from teachers" | psql

build:
	cat all.sql | psql

list:
	echo "\dt" | psql

create-table-%:
	cat $*.sql | psql
