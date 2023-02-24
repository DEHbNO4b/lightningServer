package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var queryStrikes string = `SELECT longitude,latitude,cluster FROM strikes ;`

type strike struct {
	Latitude  float32
	Longitude float32
	Cluster   int
}
type FetchStrikes struct {
	DB *sql.DB
}

func NewFetchStrikes(db *sql.DB) *FetchStrikes {
	return &FetchStrikes{db}
}
func (f FetchStrikes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	rows, err := f.DB.Query(queryStrikes)
	if err != nil {
		http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
	}
	defer rows.Close()
	var long, lat float32
	var cl int
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
	//rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	//rw.WriteHeader(200)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	rw.Write(d)
}
