import select
import sys
import psycopg2
import psycopg2.extensions

DSN="dbname='pub_sub_test' user='drio' host='localhost' password='dbpass'"
channel = "jobs_status_channel"
table = "ps_jobs"
worker_name = None

if len(sys.argv) != 2:
    print("worker name provided")
    sys.exit(1)
worker_name = sys.argv[1]

def log(msg):
    print(f'[{worker_name}] >> {msg}')

'''
FIXME: In a real world we will want to have an intermediate step
where we set the job to 'processing'
new -> processing -> success
'''
def doWork(conn, id):
    curs = conn.cursor()
    curs.execute(f'select * from {table} where id = %s;', (id,))
    rows = curs.fetchall()
    if len(rows) == 1:
        data = rows[0][1]
        status = rows[0][3]
        if status == 'new':
            addedNumbers = sum([int(c) for c in list(data) if c.isdigit()])
            log(f"Updating id={id} with result={addedNumbers}")
            curs.execute(f"""
                UPDATE
                    {table}
                SET
                    result = %s,
                    status = %s,
                    worker_name = %s,
                    status_change_time = current_timestamp
                WHERE
                    id = %s
            """, (addedNumbers, 'success', worker_name, id))
        else:
            log(f'No new job. Nothing to do.')
    curs.close()


def listen(conn):
    curs = conn.cursor()
    curs.execute(f'LISTEN {channel};')
    log(f"Waiting for notifications on channel '{channel}'")
    while True:
        if select.select([conn],[],[],5) == ([],[],[]):
            log("Timeout")
        else:
            conn.poll()
            while conn.notifies:
                notify = conn.notifies.pop(0)
                log(f"Got NOTIFY: {notify.pid}, {notify.channel}, {notify.payload}")
                doWork(conn, notify.payload)


if __name__ == '__main__':
    conn = psycopg2.connect(DSN)
    conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
    listen(conn)
