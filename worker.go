package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

type Entry struct {
	Id               int
	Data             sql.NullString
	Result           sql.NullString
	Status           sql.NullString
	WorkerName       sql.NullString
	StatusChangeTime sql.NullString
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	dsn := "postgresql://localhost:5432/pub_sub_test?sslmode=disable"
	chName := "jobs_status_channel"

	db := connect(dsn, logger)
	defer db.Close()

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logger.Println(err.Error())
		}
	}

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, reportProblem)
	err := listener.Listen(chName)
	if err != nil {
		panic(err)
	}

	logger.Println("Start monitoring PostgreSQL...")
	for {
		select {
		case n := <-listener.Notify:
			id := n.Extra
			data, status := runSelect(db, logger, id)
			if status == "new" {
				logger.Println("New Job!", id, data, status)
				doWork(db, id, data)
			}
			return
		case <-time.After(90 * time.Second):
			logger.Println("Received no events for 90 seconds, checking connection")
			go func() {
				listener.Ping()
			}()
			return
		}
	}
}

func doWork(db *sql.DB, id, data string) {
	_, err := db.Exec(`
		UPDATE
				ps_jobs
		SET
				result = $1,
				status = $2,
				worker_name = 'gopher',
				status_change_time = current_timestamp
		WHERE
				id = $3
	`, 1, "success", id)

	if err != nil {
		log.Fatalf("could not update row: %v", err)
	}
}

func runSelect(db *sql.DB, logger *log.Logger, id string) (string, string) {
	rows, err := db.Query("SELECT * from ps_jobs where id = $1", id)
	if err != nil {
		logger.Fatal(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows.Next()
	e := new(Entry)
	if err := rows.Scan(&e.Id, &e.Data, &e.Result, &e.Status, &e.WorkerName, &e.StatusChangeTime); err != nil {
		panic(err)
	}
	return e.Data.String, e.Status.String
}

func connect(dsn string, logger *log.Logger) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("database connection pool established")
	return db
}
