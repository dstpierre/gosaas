package postgres

import (
	"log"
	"os"
	"testing"

	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ds := "user=postgres password=postgres dbname=test sslmode=disable"
	conn, err := sql.Open("postgres", ds)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	// we make sure to clean everything before starting the tests
	_, err = conn.Exec("DELETE FROM gosaas_accounts;")
	if err != nil {
		log.Fatal(err)
	}

	db = conn

	retval := m.Run()
	os.Exit(retval)
}
