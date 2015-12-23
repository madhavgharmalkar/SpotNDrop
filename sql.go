package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"net/http"
	"strconv"
)

type DB struct {
	db *sql.DB
}

type Drop struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	X       float32 `json:"x"`
	Y       float32 `json:"y"`
	Date    string  `json:"date"`
	Time    string  `json:"time"`
}

func (d *DB) close() {
	d.close()
}

func (d *DB) getDropsRadius(long, lat, radx, rady float64) ([]Drop, error) {

	var minLat float64
	var maxLat float64
	var minLong float64
	var maxLong float64

	minLat = lat - (rady * 0.00904372)
	maxLat = lat + (rady * 0.00904372)

	minLong = long + radx/(math.Cos(lat*math.Pi/180)*111.320)
	maxLong = long - radx/(math.Cos(lat*math.Pi/180)*111.320)

	amin := math.Min(minLong, maxLong)
	amax := math.Max(minLong, maxLong)

	//fmt.Println(minLong)

	return d.getDropsCorners(amin, amax, minLat, maxLat)

}

func (d *DB) putDrops(dp Drop) error {
	stmt, err := d.db.Prepare("INSERT into citrus_data values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(dp.Name, dp.X, dp.Y, dp.Message, dp.Date, dp.Time)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) getDropsCorners(minLong, maxLong, minLat, maxLat float64) ([]Drop, error) {

	results := []Drop{}

	//fmt.Println("SELECT * from citrus_data where (gps_x between " + fmt.Sprint(minLong) + " and " + fmt.Sprint(maxLong) + ") and (gps_y between " + fmt.Sprint(minLat) + " and " + fmt.Sprint(maxLat) + ") order by sdate desc, stime desc;")

	rows, err := d.db.Query(
		"SELECT * from citrus_data where (gps_x between " + fmt.Sprint(minLong) + " and " + fmt.Sprint(maxLong) + ") and (gps_y between " + fmt.Sprint(minLat) + " and " + fmt.Sprint(maxLat) + ") order by sdate desc, stime desc;")
	if err != nil {
		return results, err
	}

	defer rows.Close()
	for rows.Next() {

		var a Drop
		err := rows.Scan(&a.Name, &a.X, &a.Y, &a.Message, &a.Date, &a.Time)
		if err != nil {
			return results, err
		}

		results = append(results, a)

	}

	return results, nil

}

func ConnectDB(user, password, host, port, name string) (DB, error) {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+name)
	if err != nil {
		return DB{db}, err
	}

	err = db.Ping()
	if err != nil {
		return DB{db}, err
	}

	return DB{db}, nil
}

func (d *DB) DBGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, y, radius := r.URL.Query().Get("x"), r.URL.Query().Get("y"), r.URL.Query().Get("r")

		a, err := strconv.ParseFloat(x, 64)
		if err != nil {
			handleWebErr(w, err)
			return
		}

		b, err := strconv.ParseFloat(y, 64)
		if err != nil {
			handleWebErr(w, err)
			return
		}

		c, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			handleWebErr(w, err)
			return
		}
		//fmt.Println(c)

		drops, err := d.getDropsRadius(a, b, c, c)
		if err != nil {
			handleWebErr(w, err)
			return
		}

		d, err := json.Marshal(drops)
		if err != nil {
			handleWebErr(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(d))

	}
}
