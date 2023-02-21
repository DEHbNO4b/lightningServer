package main

import (
	"database/sql"
	"lServer/handlers"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/stdlib"
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

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	s.ListenAndServe()

	// tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// s.Shutdown(tc)
}
