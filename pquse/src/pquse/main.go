package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:password@192.168.1.50:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM weather")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			city string
			lo   int
			hi   int
			prcp float64
			date string
		)

		err = rows.Scan(&city, &lo, &hi, &prcp, &date)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("\ncity %s\nlo   %d\nhi   %d\nprcp %f\ndate %s\n", city, lo, hi, prcp, date)
	}

	log.Printf("over\n")
}
