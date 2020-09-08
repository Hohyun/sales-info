package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "hohyunkim"
	password = "Django121*"
	dbname   = "selabd"
)

var (
	salesdate   string
	domintl     string
	salesrefund string
	krw_amt     float64
)

// func ImportCsv(src_file string) {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal("Failed to open a DB connection: ", err)
// 	}
// 	defer db.Close()

// 	_, err := db.Exec(fmt.Sprintf("COPY sales FROM %s DELIMITER ',' CSV HEADER", src_file))
// 	if err != nil {
// 		log.Fatal("Failed to import csv file: ", err)
// 	}
// }

func QuerySales() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	fmt.Println("* Sales by DOM/INTL")
	rows, err := db.Query("select domintl, sum(krwamt) krw_amount from sales group by domintl order by domintl")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&domintl, &krw_amt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s --> %13.0f\n", domintl, krw_amt)
	}

	fmt.Println("\n* Sales by DOM/INTL, Sales/Refund")
	rows, err = db.Query("select domintl, salesrefund, sum(krwamt) krw_amount from sales group by domintl, salesrefund order by domintl, salesrefund desc")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&domintl, &salesrefund, &krw_amt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %-10s --> %13.0f\n", domintl, salesrefund, krw_amt)
	}

	fmt.Println("\n* Sales by Date, DOM/INTL, Sales/Refund")
	rows, err = db.Query("select salesdate, domintl, salesrefund, sum(krwamt) krw_amount from sales group by salesdate, domintl, salesrefund order by salesdate, domintl, salesrefund desc")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&salesdate, &domintl, &salesrefund, &krw_amt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %-10s %-10s --> %13.0f\n", salesdate[0:10], domintl, salesrefund, krw_amt)
	}
}
