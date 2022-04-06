package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

/* -- this varibles defined in query-postgres.go
var (
	salesdate   string
	agencytype  string
	fop         string
	fopdesc     string
	domintl     string
	ccy         string
	salesrefund string
	isales      float64
	irfnd       float64
	iyr         float64
	itax        float64
	itotal      float64
	dsales      float64
	drfnd       float64
	dyr         float64
	dtax        float64
	dtotal      float64
	gtotal      float64
	krwamt      float64
	amount      float64
) */

// ImportCsvSqSales : copy csv file into database.
func ImportCsvSqSales(cfg Config) {
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	r1, err := db.Exec("DELETE FROM sales_raw")
	if err != nil {
		log.Fatal(err)
	}
	affected1, _ := r1.RowsAffected()
	fmt.Printf("Deleted  : %d rows\n", affected1)

	// cmd := exec.Command("sqlite3", "./data/selabd.db", ".read import.sql")
	cmd := exec.Command("./sqlite3.exe", cfg.SqliteDb, ".read import_sales.sql")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to import to sales_raw table ", err)
	}

	sqlStr := `
	INSERT INTO sales
    SELECT salesdate, agencytype, fop, fopdesc, domintl, salesrefund, ccy, sum(amount) amount, sum(krwamt) krwamt
    FROM sales_raw GROUP BY 1,2,3,4,5,6,7
    ON CONFLICT (salesdate, agencytype, fop, domintl, salesrefund, ccy)
    DO UPDATE SET amount = excluded.amount, krwamt = excluded.krwamt;`
	affected2, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal("Failed to insert summarized data into sales table: ", err)
	}
	fmt.Printf("Inserted to sales table: %v rows\n", affected2)
	fmt.Printf("Data import finished successfully\n")
}

// ImportCsvSqTaxYr : copy csv file into database.
func ImportCsvSqTaxYr(cfg Config) {
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	r1, err := db.Exec("DELETE FROM taxyr_raw")
	if err != nil {
		log.Fatal(err)
	}
	affected1, _ := r1.RowsAffected()
	fmt.Printf("Deleted  : %d rows\n", affected1)

	cmd := exec.Command("./sqlite3.exe", cfg.SqliteDb, ".read import_taxyr.sql")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to import to taxyr_raw table ", err)
	}

	sqlStr := `
	INSERT INTO taxyr
	SELECT salesdate, taxyr, domintl, ccy, sum(salesamt) salesamt, sum(refundamt) refundamt, sum(reissueamt) reissueamt
	FROM taxyr_raw
	GROUP BY salesdate, taxyr, domintl, ccy
	ORDER BY salesdate, taxyr, domintl, ccy
	ON CONFLICT (salesdate, taxyr, domintl, ccy)
	DO UPDATE set salesamt = excluded.salesamt, refundamt = excluded.refundamt, reissueamt = excluded.reissueamt;`
	affected2, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal("Failed to insert summarized data into taxyr table: ", err)
	}
	fmt.Printf("Inserted to taxyr table: %v rows\n", affected2)
	fmt.Printf("Data import finished successfully\n")
}