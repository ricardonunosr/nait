package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	supa "github.com/nedpals/supabase-go"
)

var (
	Db             *sql.DB
	StaffList      *sql.Stmt
	EventsList     *sql.Stmt
	EventsNameList *sql.Stmt
)

func InitDB() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(256)
	db.SetMaxOpenConns(256)

	Db = db

	StaffList, err = db.Prepare("SELECT * FROM staff;")
	if err != nil {
		panic(err)
	}
	EventsList, err = db.Prepare("SELECT e.event_date, en.event_name FROM events e INNER JOIN events_name en ON e.event_id = en.event_id")
	if err != nil {
		panic(err)
	}

	EventsNameList, err = db.Prepare("SELECT * FROM events_name")
	if err != nil {
		panic(err)
	}
}

func CreateSupabaseClient() *supa.Client {
	return supa.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}
