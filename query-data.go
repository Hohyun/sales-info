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
	krwamt      float64
)

// ImportCsv copy csv file into database.
func ImportCsv(srcFile string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

func sales(db *sql.DB) {
	row := db.QueryRow("select sum(krwamt) krw_amount from sales")
	fmt.Printf("%-7s %-10s %-4s %-6s %13s\n", "Level 0", "", "D/I", "", "KRW Amount")
	fmt.Println("------- ---------- ---- ------ -------------")
	err := row.Scan(&krwamt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-7s %-10s %-4s %-6s %13.0f\n", "", "", "", "", krwamt)
	fmt.Printf("\n")
}

func salesDomIntl(db *sql.DB) {
	rows, err := db.Query("select domintl, sum(krwamt) krw_amount from sales group by domintl order by domintl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-7s %-10s %-4s %-6s %13s\n", "Level 1", "", "D/I", "", "KRW Amount")
	fmt.Println("------- ---------- ---- ------ -------------")
	for rows.Next() {
		err := rows.Scan(&domintl, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-7s %-10s %-4s %-6s %13.0f\n", "", "", domintl, "", krwamt)
	}
	fmt.Printf("\n")
}

func salesDomIntlSalesRefund(db *sql.DB) {
	rows, err := db.Query("select domintl, salesrefund, sum(krwamt) krw_amount from sales group by domintl, salesrefund order by domintl, salesrefund desc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-7s %-10s %-4s %-6s %13s\n", "Level 2", "", "D/I", "S/R", "KRW Amount")
	fmt.Println("------- ---------- ---- ------ -------------")
	for rows.Next() {
		err := rows.Scan(&domintl, &salesrefund, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-7s %-10s %-4s %-6s %13.0f\n", "", "", domintl, salesrefund, krwamt)
	}
	fmt.Printf("\n")
}

func salesDateDomIntlSalesRefund(db *sql.DB) {
	rows, err := db.Query("select salesdate, domintl, salesrefund, sum(krwamt) krw_amount from sales group by salesdate, domintl, salesrefund order by salesdate, domintl, salesrefund desc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-7s %-10s %-4s %-6s %13s\n", "Level 3", "Sales Date", "D/I", "S/R", "KRW Amount")
	fmt.Println("------- ---------- ---- ------ -------------")
	for rows.Next() {
		err := rows.Scan(&salesdate, &domintl, &salesrefund, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-7s %-10s %-4s %-6s %13.0f\n", "", salesdate[0:10], domintl, salesrefund, krwamt)
	}
	fmt.Printf("\n")
}

// QuerySales show query results from database
func QuerySales(level int) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch level {
	case 0:
		sales(db)
	case 1:
		salesDomIntl(db)
	case 2:
		salesDomIntlSalesRefund(db)
	case 3:
		salesDateDomIntlSalesRefund(db)
	default:
		sales(db)
		salesDomIntl(db)
		salesDomIntlSalesRefund(db)
		salesDateDomIntlSalesRefund(db)
	}
}
