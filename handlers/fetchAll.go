package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

var queryStrikes string = `SELECT longitude,latitude,cluster FROM strikes;`

type strike struct {
	latitude  float32
	longitude float32
	cluster   int
}
type FetchAll struct {
	DB *sql.DB
}

func NewFetchAll(db *sql.DB) *FetchAll {
	return &FetchAll{db}
}
func (f FetchAll) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rows, err := f.DB.Query(queryStrikes)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var long, lat float32
	var cl int
	var data []strike
	for rows.Next() {
		if err = rows.Scan(&long, &lat, &cl); err != nil {
			panic(err)
		}
		s := strike{longitude: long, latitude: lat, cluster: cl}
		data = append(data, s)
	}
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}
