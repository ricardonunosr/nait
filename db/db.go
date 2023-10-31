package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	supa "github.com/nedpals/supabase-go"
)

var (
	DB               *sql.DB
	STAFF_LIST       *sql.Stmt
	EVENTS_LIST      *sql.Stmt
	EVENTS_NAME_LIST *sql.Stmt
)

func InitDB() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(256)
	db.SetMaxOpenConns(256)

	DB = db

	STAFF_LIST, err = db.Prepare("SELECT * FROM staff;")
	if err != nil {
		panic(err)
	}
	EVENTS_LIST, err = db.Prepare("SELECT e.event_date, en.event_name FROM events e INNER JOIN events_name en ON e.event_id = en.event_id")
	if err != nil {
		panic(err)
	}

	EVENTS_NAME_LIST, err = db.Prepare("SELECT * FROM events_name")
	if err != nil {
		panic(err)
	}
}

func CreateSupabaseClient() *supa.Client {
	return supa.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}
