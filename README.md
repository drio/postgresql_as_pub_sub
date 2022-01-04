## Intro

This is a proof of concept to explore the idea of using [PostgreSQL as a job server](https://webapp.io/blog/postgres-is-the-answer/).
server.

<img width="1904" alt="Screen Shot 2022-01-04 at 8 52 11 AM" src="https://user-images.githubusercontent.com/17954/148070510-e153b0cb-3eba-44d1-ae29-afea90da2736.png">

## Usage

1. Start your psql instance: `$ make start-fg` (console 1)
2. Create the database and tables: `$ make clean up-sql`
3. Create jobs: `$ while :; do make new_job list_jobs; sleep 5; done` (console 2)
4. Connect to the db with the psql client: `make` (console 3) and tell psql you 
   want to listen to the channel that omits events on table changes with: `listen jobs_status_channel;`.
   This is just for you to make sure you are getting data in the channel.
6. In another console, run the woker (console 4):
    ```
    $ source venv/bin/activate
    $ python ./worker.py
    ```

### Home brew installation output

```txt
==> postgresql
To migrate existing data from a previous major version of PostgreSQL run:
  brew postgresql-upgrade-database

This formula has created a default database cluster with:
  initdb --locale=C -E UTF-8 /opt/homebrew/var/postgres
For more details, read:
  https://www.postgresql.org/docs/14/app-initdb.html

To restart postgresql after an upgrade:
  brew services restart postgresql
Or, if you don't want/need a background service you can just run:
  /opt/homebrew/opt/postgresql/bin/postgres -D /opt/homebrew/var/postgres
```

## References

1. [Psql as job server](https://webapp.io/blog/postgres-is-the-answer/)
2. [HN discussion](https://news.ycombinator.com/item?id=29599132)
3. [Psql cheatsheet](https://gist.github.com/xpepper/8110743)