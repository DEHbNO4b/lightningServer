package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type thunder struct {
	Id        int
	Claster   string
	Polygon   [][]float32
	Area      float32   `json:"area"`
	Capacity  int       `json:"capacity"`
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
	Duration  time.Duration
}

type FetchThunders struct {
	DB *sql.DB
}

func NewFetchThunders(db *sql.DB) FetchThunders {
	return FetchThunders{db}
}

func (f FetchThunders) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	requestedDay := r.URL.Query().Get("day")
	day, err := time.Parse("2006.01.02_MST", requestedDay)
	if err != nil {
		http.Error(rw, "Unable purse requested day", http.StatusInternalServerError)
	}
	day2 := day.AddDate(0, 0, 1)
	var id, capacity int
	var area float32
	var startTime, endTime time.Time
	var claster, p string
	var data []thunder

	rows, err := f.DB.Query(`SELECT id,cluster,st_asText(geog),area,capacity,startTime,endTime 
							FROM thunders WHERE starttime BETWEEN $1 and $2;`, day, day2)
	if err != nil {
		http.Error(rw, "Unable to connect DB", http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &claster, &p, &area, &capacity, &startTime, &endTime); err != nil {
			http.Error(rw, "Unable read data from db", http.StatusInternalServerError)
		}
		points := parsePolygon(p)
		t := thunder{Id: id, Claster: claster, Polygon: points, Area: area, Capacity: capacity, StartTime: startTime, EndTime: endTime}
		data = append(data, t)
	}
	d, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	rw.Write(d)
}

func parsePolygon(p string) [][]float32 {
	var points [][]float32
	p, _ = strings.CutPrefix(p, "POLYGON((")
	p, _ = strings.CutSuffix(p, "))")
	poly := strings.Split(p, ",")
	for _, val := range poly {
		s := strings.Fields(val)

		long, err := strconv.ParseFloat(s[0], 32)
		if err != nil {
			fmt.Println("Unable convert coordinates")
		}
		lat, err := strconv.ParseFloat(s[1], 32)
		if err != nil {
			fmt.Println("Unable convert coordinates")
		}
		points = append(points, []float32{float32(lat), float32(long)})
	}
	return points
}
