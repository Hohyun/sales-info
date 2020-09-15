package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

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
)

// ImportCsvPG : copy csv file into database.
func ImportCsvPG(srcFile string) {
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
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

	r2, err := db.Exec(fmt.Sprintf("COPY sales FROM '%s' DELIMITER ',' CSV HEADER", srcFile))
	if err != nil {
		log.Fatal("Failed to import csv file: ", err)
	}
	affected2, _ := r2.RowsAffected()
	fmt.Printf("Inserted: %d rows\n", affected2)
}

// ExportCsvPG : export sales results into csv file.
func ExportCsvPG(dstFile string) {
	// Open database
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
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

func salesRawPG(db *sql.DB) {
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

func salesTabularPG(db *sql.DB) {
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

// QuerySalesPG show query results from database
func QuerySalesPG(reportType string) {
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch reportType {
	case "table":
		salesTabularPG(db)
	case "raw":
		salesRawPG(db)
	default:
		salesTabularPG(db)
	}
}
