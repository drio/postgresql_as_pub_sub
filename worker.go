package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

/*
	1. Listen for channel events
	2. Use the id from channel and select the table to get new job
	3. do the work and perform an update
*/
func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//db, err := sql.Open("postgres", "user=drio password= dbname=pub_sub_test sslmode=disable")
	dsn := "postgresql://localhost:5432/pub_sub_test?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	entries := make([]*Entry, 0)
	rows, err := db.Query("SELECT * from ps_jobs")
	if err != nil {
		logger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		e := new(Entry)
		if err := rows.Scan(&e.Id, &e.Data, &e.Result, &e.Status, &e.WorkerName, &e.StatusChangeTime); err != nil {
			panic(err)
		}
		entries = append(entries, e)
		fmt.Println(e.Id, e.Data.String)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	/* Test LISTEN */
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, reportProblem)
	chName := "jobs_status_channel"
	err = listener.Listen(chName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	waitForNotification(listener)
	fmt.Println("Done")
}

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			// Prepare notification payload for pretty print
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			id := string(prettyJSON.Bytes())
			fmt.Println("Id is : ", id)
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
