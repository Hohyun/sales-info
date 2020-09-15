package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/exec"
	// "path"
	// "strings"
)

/* -- this varibles defined in query-postgres.go
var (
	salesdate   string
	isales      float64
	irfnd       float64
	itotal      float64
	dsales      float64
	drfnd       float64
	dtotal      float64
	gtotal      float64
	domintl     string
	salesrefund string
	krwamt      float64
) */

// ImportCsvSQ : copy csv file into database.
func ImportCsvSQ(srcFile string, cfg Config) {
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	r1, err := db.Exec("DELETE FROM sales")
	if err != nil {
		log.Fatal("Failed to import csv file: ", err)
	}
	affected1, _ := r1.RowsAffected()
	fmt.Printf("Deleted  : %d rows\n", affected1)

	// .import
	if err != nil {
		log.Fatal(err)
	}
	// cmd := exec.Command("sqlite3", "./data/selabd.db", ".read import.sql")
	cmd := exec.Command("./sqlite3.exe", cfg.SqliteDb, ".read import.sql")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data imported successfully\n")
}

// ExportCsvSQ : export sales results into csv file.
func ExportCsvSQ(dstFile string, cfg Config) {
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	// Get data: Sale by date
	rows, err := db.Query("select * from sales_by_date")
	if err != nil {
		panic(err)
	}

	var ss [][]string
	ss = append(ss, []string{"Date", "Intl_sales", "Intl_refund", "Intl_total", "Dom_sales", "Dom_refund", "Dom_total", "G_total"})
	for rows.Next() {
		var s []string
		err := rows.Scan(&salesdate, &isales, &irfnd, &itotal, &dsales, &drfnd, &dtotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, salesdate[0:10], fmt.Sprintf("%.0f", isales), fmt.Sprintf("%.0f", irfnd),
			fmt.Sprintf("%.0f", itotal), fmt.Sprintf("%.0f", dsales), fmt.Sprintf("%.0f", drfnd),
			fmt.Sprintf("%.0f", dtotal), fmt.Sprintf("%.0f", gtotal))
		ss = append(ss, s)
	}

	// Get data: Summary line
	row := db.QueryRow("select * from sales_summary")
	err = row.Scan(&isales, &irfnd, &itotal, &dsales, &drfnd, &dtotal, &gtotal)
	if err != nil {
		panic(err)
	}
	var s []string
	s = append(s, "Total", fmt.Sprintf("%.0f", isales), fmt.Sprintf("%.0f", irfnd),
		fmt.Sprintf("%.0f", itotal), fmt.Sprintf("%.0f", dsales), fmt.Sprintf("%.0f", drfnd),
		fmt.Sprintf("%.0f", dtotal), fmt.Sprintf("%.0f", gtotal))
	ss = append(ss, s)

	// Write csv file
	f, err := os.Create(dstFile)
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	w := csv.NewWriter(f)
	defer f.Close()
	err = w.WriteAll(ss)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sales results was exported to %s successfully!\n", dstFile)
}

func salesRawSQ(db *sql.DB) {
	rows, err := db.Query("select salesdate, domintl, salesrefund, sum(krwamt) krw_amount from sales group by salesdate, domintl, salesrefund order by salesdate, domintl, salesrefund desc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-10s %-4s %-6s %12s\n", "Date", "D/I", "S/R", "KRW Amount")
	fmt.Println("---------- ---- ------ -------------")
	for rows.Next() {
		err := rows.Scan(&salesdate, &domintl, &salesrefund, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %-4s %-6s %12.0f\n", salesdate[0:10], domintl, salesrefund, krwamt)
	}
	fmt.Printf("\n")
}

func salesTabularSQ(db *sql.DB) {
	// Summary by date
	rows, err := db.Query("select * from sales_by_date")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n-----------------------------------------------------------------------------------------------------\n")
	fmt.Printf("%-10s %38s %38s\n", "", "INTL", "DOM")
	fmt.Printf("           -------------------------------------- --------------------------------------\n")
	fmt.Printf("%-10s %12s %12s %12s %12s %12s %12s %12s\n",
		"Date", "Sales", "Refund", "Total", "Sales", "Refund", "Total", "G.Total")
	fmt.Printf("---------- ------------ ------------ ------------ ------------ ------------ ------------ ------------\n")
	for rows.Next() {
		err := rows.Scan(&salesdate, &isales, &irfnd, &itotal, &dsales, &drfnd, &dtotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %12.0f %12.0f %12.0f %12.0f %12.0f %12.0f %12.0f\n",
			salesdate[0:10], isales, irfnd, itotal, dsales, drfnd, dtotal, gtotal)
	}

	// Total summary line
	row := db.QueryRow("select * from sales_summary")
	err = row.Scan(&isales, &irfnd, &itotal, &dsales, &drfnd, &dtotal, &gtotal)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("---------- ------------ ------------ ------------ ------------ ------------ ------------ ------------\n")
	fmt.Printf("%-10s %12.0f %12.0f %12.0f %12.0f %12.0f %12.0f %12.0f\n",
		"Total", isales, irfnd, itotal, dsales, drfnd, dtotal, gtotal)
	fmt.Printf("---------- ------------ ------------ ------------ ------------ ------------ ------------ ------------\n")
	fmt.Printf("\n")
}

// QuerySalesSQ show query results from database
func QuerySalesSQ(reportType string, cfg Config) {
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch reportType {
	case "table":
		salesTabularSQ(db)
	case "raw":
		salesRawSQ(db)
	default:
		salesTabularSQ(db)
	}
}
