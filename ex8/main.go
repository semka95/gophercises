package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./sqlite.db")
	defer db.Close()
	createTable := `CREATE TABLE IF NOT EXISTS phone (
"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
"phone" TEXT NOT NULL UNIQUE
);`
	st, err := db.Prepare(createTable)
	if err != nil {
		panic(err)
	}
	_, err = st.Exec()
	if err != nil {
		panic(err)
	}
}
