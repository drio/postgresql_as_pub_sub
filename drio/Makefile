var/postgres:
	initdb ./var/postgres

start-fg: var/postgres
	/opt/homebrew/opt/postgresql/bin/postgres -D ./var/postgres

start-fg-default:
	/opt/homebrew/opt/postgresql/bin/postg es -D /opt/homebrew/var/postgres
