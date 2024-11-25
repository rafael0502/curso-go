package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Conectar ao banco de dados SQLite
	db, err := sql.Open("sqlite3", "./teste1.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	// Criar uma tabela
	sqlStmt := ` CREATE TABLE IF NOT EXISTS users ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT );`

	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	stmt, err := db.Prepare("INSERT INTO users(name) VALUES (?)")

	if err != nil {
		log.Fatalf("%q\n", err)
	}

	defer stmt.Close()
	_, err = stmt.Exec("Rafa")

	if err != nil {
		log.Fatalf("%q\n", err)
	}

	log.Println("Fim")
}
