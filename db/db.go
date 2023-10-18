package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	supa "github.com/nedpals/supabase-go"
)

var DB *sql.DB

func InitDB() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err.Error())
	}

	DB = db
}

func CreateSupabaseClient() *supa.Client {
	return supa.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}
