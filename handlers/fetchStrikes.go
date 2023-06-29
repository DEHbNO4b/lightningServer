package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type strike struct {
	Latitude  float32
	Longitude float32
	Cluster   string
}
type FetchStrikes struct {
	DB *sql.DB
}

func NewFetchStrikes(db *sql.DB) *FetchStrikes {
	return &FetchStrikes{db}
}
func (f FetchStrikes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	requestedDay := r.URL.Query().Get("day")
	day, err := time.Parse("2006.01.02_MST", requestedDay)
	if err != nil {
		http.Error(rw, "Unable purse requested day", http.StatusInternalServerError)
	}
	//day.Add(-3)
	day2 := day.AddDate(0, 0, 1)
	fmt.Println("day;", day)
	fmt.Println("date2;", day2)

	rows, err := f.DB.Query(`SELECT longitude,latitude,cluster FROM strikes where time between $1 and $2`, day, day2)
	if err != nil {
		http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
	}
	defer rows.Close()
	var long, lat float32
	var cl string
	var data []strike
	for rows.Next() {
		if err = rows.Scan(&long, &lat, &cl); err != nil {
			http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
		}
		s := strike{Longitude: long, Latitude: lat, Cluster: cl}

		data = append(data, s)
	}
	d, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	rw.Write(d)
}
