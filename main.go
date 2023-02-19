package main

import (
	"database/sql"
	"lServer/handlers"
	"net/http"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"

func main() {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fa := handlers.NewFetchAll(db)
	sm := http.NewServeMux()
	sm.Handle("/", fa)
	http.ListenAndServe(":9090", sm)
}
