package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/goneri/go_1w"
	_ "github.com/lib/pq"
)

func init_db(db *sql.DB) {
	create_table := "CREATE TABLE IF NOT EXISTS temps (id serial PRIMARY KEY, value FLOAT(6), created_at  TIMESTAMP NOT NULL DEFAULT NOW());"
	_, err := db.Exec(create_table)
	checkErr(err)
	//fmt.Println(result)

}

func recorder(db *sql.DB) {

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				record(db, go_1w.Read("sample"))
				// case <-quit:
				//	ticker.Stop()
				//	return
			}
		}
	}()

}

func record(db *sql.DB, value float32) {
	_, err := db.Exec(`INSERT into temps (value) values ($1)`, value)
	checkErr(err)
}

func get_temp(db *sql.DB) float32 {
	row := db.QueryRow(`select value from temps order by created_at desc limit 1;`)
	var value float32
	err := row.Scan(&value)
	checkErr(err)
	return value
}

func handler_index(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("%f", get_temp(db)))
	}
}

func main() {
	db, err := sql.Open("postgres", "user=goneri host=/run/postgresql dbname=goneri sslmode=disable")
	checkErr(err)
	defer db.Close()
	init_db(db)
	recorder(db)

	http.HandleFunc("/", handler_index(db))

	http.ListenAndServe(":8080", nil)
	//time.Sleep(16000 * time.Millisecond)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
