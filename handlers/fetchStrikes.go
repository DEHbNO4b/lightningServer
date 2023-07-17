package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type strike struct {
	Latitude  float32
	Longitude float32
	Cluster   string
	Id        int
}
type FetchStrikes struct {
	DB *sql.DB
}

func NewFetchStrikes(db *sql.DB) *FetchStrikes {
	return &FetchStrikes{db}
}
func (f FetchStrikes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	requestedDay := r.URL.Query().Get("day")
	day, err := time.Parse("2006.01.02_15:04_MST", requestedDay)
	if err != nil {
		http.Error(rw, "Unable purse requested day", http.StatusInternalServerError)
	}
	// day2 := day.AddDate(0, 0, 1)
	time2 := day.Add(10 * time.Minute)
	rows, err := f.DB.Query(`SELECT longitude,latitude,cluster,id FROM enstrikes where time >= $1 and time<$2`, day, time2)
	if err != nil {
		http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
	}
	defer rows.Close()
	var long, lat float32
	var id int
	var cl string
	var data []strike
	for rows.Next() {
		if err := rows.Scan(&long, &lat, &cl, &id); err != nil {

			// http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
			fmt.Println(err)
			// return
		}
		s := strike{Longitude: long, Latitude: lat, Cluster: cl, Id: id}
		fmt.Println(cl)
		if len(s.Cluster) > 2 {
			c := s.Cluster[len(s.Cluster)-2:]
			c, _ = strings.CutPrefix(c, ":")
			s.Cluster = c
		}

		data = append(data, s)
	}
	fmt.Println("count: ", len(data))
	d, err := json.Marshal(data)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	rw.WriteHeader(http.StatusOK)
	rw.Write(d)
}
