package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FetchDayCollection struct {
	DB *sql.DB
}

func NewFetchDayCollection(db *sql.DB) *FetchDayCollection {
	return &FetchDayCollection{db}
}
func (f FetchDayCollection) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var dayCollection []string
	rows, err := f.DB.Query(`select distinct date_trunc('day',time) from strikes ORDER BY date_trunc;`)
	if err != nil {
		http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
	}
	defer rows.Close()
	var dateOnly time.Time

	for rows.Next() {
		if err = rows.Scan(&dateOnly); err != nil {
			http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
		}
		fmt.Println(dateOnly)
		dayCollection = append(dayCollection, dateOnly.Format("2006.01.02_MST"))
	}

	data, err := json.Marshal(dayCollection)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	rw.Write(data)
}
