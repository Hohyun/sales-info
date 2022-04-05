package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	ON CONFLICT (salesdate, taxyr, domintl, ccy)
	DO UPDATE set salesamt = excluded.salesamt, refundamt = excluded.refundamt, reissueamt = excluded.reissueamt;`
	affected2, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal("Failed to insert summarized data into taxyr table: ", err)
	}
	fmt.Printf("Inserted to taxyr table: %v rows\n", affected2)
	fmt.Printf("Data import finished successfully\n")
}

func exportRawSQ(db *sql.DB, dstFile string, fromDate string, toDate string) {
	sqlStr := fmt.Sprintf(`select salesdate, agencytype, fop, fopdesc, domintl, salesrefund, ccy, amount, krwamt
	from sales where salesdate between '%s' and '%s'`, fromDate, toDate)
	rows, err := db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	var ss [][]string
	ss = append(ss, []string{"SalesDate", "AgencyType", "Fop", "FopDesc", "DomIntl", "S_R", "Ccy", "Amount", "KrwAmount"})
	for rows.Next() {
		var s []string
		err := rows.Scan(&salesdate, &agencytype, &fop, &fopdesc, &domintl, &salesrefund, &ccy, &amount, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, salesdate[0:10], agencytype, fop, fopdesc, domintl, salesrefund, ccy,
			fmt.Sprintf("%.0f", amount), fmt.Sprintf("%.0f", krwamt))
		ss = append(ss, s)
	}
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

func exportTabularSQ1(db *sql.DB, dstFile string, fromDate string, toDate string) {
	// Get data: Sale by date
	rows, err := db.Query(fmt.Sprintf("select * from sales_tax_yr_by_date where salesdate between '%s' and '%s'", fromDate, toDate))
	if err != nil {
		panic(err)
	}

	var ss [][]string
	var dsalesT, dyrT, dtaxT, dtotalT, isalesT, iyrT, itaxT, itotalT, gtotalT float64
	ss = append(ss, []string{"Date", "Sales", "YR Rev", "Tax", "Total", "Sales", "YR Rev", "Tax", "Total", "G.Total"})
	
	for rows.Next() {
		var s1 []string
		err := rows.Scan(&salesdate, &dsales, &dyr, &dtax, &dtotal, &isales, &iyr, &itax, &itotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		s1 = append(s1, salesdate[0:10], fmt.Sprintf("%.0f", dsales), fmt.Sprintf("%.0f", dyr), fmt.Sprintf("%.0f", dtax),
			fmt.Sprintf("%.0f", dtotal), fmt.Sprintf("%.0f", isales), fmt.Sprintf("%.0f", iyr), fmt.Sprintf("%.0f", itax),
			fmt.Sprintf("%.0f", itotal), fmt.Sprintf("%.0f", gtotal))
		ss = append(ss, s1)
		dsalesT += dsales
		dyrT += dyr
		dtaxT += dtax
		dtotalT += dtotal
		isalesT += isales
		iyrT += iyr
		itaxT += itax
		itotalT += itotal
		gtotalT += gtotal
	}
	var s2 []string
	s2 = append(s2, "Total", fmt.Sprintf("%.0f", dsalesT), fmt.Sprintf("%.0f", dyrT), fmt.Sprintf("%.0f", dtaxT),
		fmt.Sprintf("%.0f", dtotalT), fmt.Sprintf("%.0f", isalesT), fmt.Sprintf("%.0f", iyrT), fmt.Sprintf("%.0f", itaxT),
		fmt.Sprintf("%.0f", itotalT), fmt.Sprintf("%.0f", gtotalT))
	ss = append(ss, s2)

	// Write csv file
	dstFile = strings.Replace(dstFile, ".", "_sales_yr_tax.", 1)
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

func exportTabularSQ2(db *sql.DB, dstFile string, fromDate string, toDate string) {
	// Get data: Sale by date
	rows, err := db.Query(fmt.Sprintf("select * from sales_by_date where salesdate between '%s' and '%s'", fromDate, toDate))
	if err != nil {
		panic(err)
	}

	var ss [][]string
	var dsalesT, drfndT, isalesT, irfndT float64
	ss = append(ss, []string{"Date", "Dom_sales", "Dom_refund", "Dom_total", "Intl_sales", "Intl_refund", "Intl_total", "G_total"})
	for rows.Next() {
		var s1 []string
		err := rows.Scan(&salesdate, &dsales, &drfnd, &dtotal, &isales, &irfnd, &itotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		s1 = append(s1, salesdate[0:10], fmt.Sprintf("%.0f", dsales), fmt.Sprintf("%.0f", drfnd),
			fmt.Sprintf("%.0f", dtotal), fmt.Sprintf("%.0f", isales), fmt.Sprintf("%.0f", irfnd),
			fmt.Sprintf("%.0f", itotal), fmt.Sprintf("%.0f", gtotal))
		ss = append(ss, s1)
		dsalesT += dsales
		drfndT += drfnd
		isalesT += isales
		irfndT += irfnd
	}
	var s2 []string
	s2 = append(s2, "Total", fmt.Sprintf("%.0f", dsalesT), fmt.Sprintf("%.0f", drfndT),
		fmt.Sprintf("%.0f", dsalesT+drfndT), fmt.Sprintf("%.0f", isalesT), fmt.Sprintf("%.0f", irfndT),
		fmt.Sprintf("%.0f", isalesT+irfndT), fmt.Sprintf("%.0f", dsalesT+drfndT+isalesT+irfndT))
	ss = append(ss, s2)

	// Write csv file
	dstFile = strings.Replace(dstFile, ".", "_sales_rfnd.", 1)
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

// ExportCsvSQ : export sales data with different format.
func ExportCsvSQ(reportType string, dstFile string, fromDate string, toDate string, cfg Config) {
	// Open database
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch reportType {
	case "tabular":
		exportTabularSQ1(db, dstFile, fromDate, toDate)
		exportTabularSQ2(db, dstFile, fromDate, toDate)
	case "raw":
		exportRawSQ(db, dstFile, fromDate, toDate)
	default:
		exportTabularSQ1(db, dstFile, fromDate, toDate)
		exportTabularSQ2(db, dstFile, fromDate, toDate)
	}
}

func salesRawSQ(db *sql.DB, fromDate string, toDate string) {
	sqlStr := fmt.Sprintf(`select salesdate, agencytype, fop, fopdesc, domintl, salesrefund, ccy, amount, krwamt
	from sales where salesdate between '%s' and '%s'`, fromDate, toDate)
	rows, err := db.Query(sqlStr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-10s %-14s %-10s %-30s %-4s %-6s %3s %14s %12s\n",
		"Date", "AgencyType", "FOP", "FOP Desc", "DomIntl", "S_R", "Ccy", "Amount", "KrwAmount")
	fmt.Println("---------- -------------- ---------- -------------------- ---- ------ --- -------------- ------------")
	for rows.Next() {
		err := rows.Scan(&salesdate, &agencytype, &fop, &fopdesc, &domintl, &salesrefund, &ccy, &amount, &krwamt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %-14s %-10s %-30s %-4s %-6s %3s %14.2f %12.0f\n",
			salesdate[0:10], agencytype, fop, fopdesc, domintl, salesrefund, ccy, amount, krwamt)
	}
	fmt.Printf("\n")
}

func salesTabularSQ1(db *sql.DB, fromDate string, toDate string) {
	var dsalesT, dyrT, dtaxT, dtotalT, isalesT, iyrT, itaxT, itotalT, gtotalT float64

	rows, err := db.Query(fmt.Sprintf("select * from sales_tax_yr_by_date where salesdate between '%s' and '%s'", fromDate, toDate))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n-------------------------------------------------------------------------------------------------------------------------------------------------\n")
	fmt.Printf("%-10s %59s %59s\n", "", "DOM", "INTL")
	fmt.Printf("           ---------------------------------------------------------- ---------------------------------------------------------- \n")
	fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s %14s %14s\n",
		"Date", "Sales", "YR Rev", "Tax", "Total", "Sales", "YR Rev", "Tax", "Total", "G.Total")
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	for rows.Next() {
		err := rows.Scan(&salesdate, &dsales, &dyr, &dtax, &dtotal, &isales, &iyr, &itax, &itotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s %14s %14s\n",
			salesdate[0:10], commas(int(dsales)), commas(int(dyr)), commas(int(dtax)), commas(int(dtotal)), 
			commas(int(isales)), commas(int(iyr)), commas(int(itax)), commas(int(itotal)), commas(int(gtotal)))
		dsalesT += dsales
		dyrT += dyr
		dtaxT += dtax
		dtotalT += dtotal
		isalesT += isales
		iyrT += iyr
		itaxT += itax
		itotalT += itotal
		gtotalT += gtotal
	}
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s %14s %14s\n",
		"Total", commas(int(dsalesT)), commas(int(dyrT)), commas(int(dtaxT)), commas(int(dtotalT)), 
		commas(int(isalesT)), commas(int(iyrT)), commas(int(itaxT)), commas(int(itotalT)), commas(int(gtotalT)))
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	fmt.Printf("\n")
}

func salesTabularSQ2(db *sql.DB, fromDate string, toDate string) {
	var dsalesT, drfndT, isalesT, irfndT float64

	rows, err := db.Query(fmt.Sprintf("select * from sales_by_date where salesdate between '%s' and '%s'", fromDate, toDate))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n-------------------------------------------------------------------------------------------------------------------\n")
	fmt.Printf("%-10s %44s %44s\n", "", "DOM", "INTL")
	fmt.Printf("           -------------------------------------------- --------------------------------------------\n")
	fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s\n",
		"Date", "Sales", "Refund", "Total", "Sales", "Refund", "Total", "G.Total")
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	for rows.Next() {
		err := rows.Scan(&salesdate, &dsales, &drfnd, &dtotal, &isales, &irfnd, &itotal, &gtotal)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s\n",
			salesdate[0:10], commas(int(dsales)), commas(int(drfnd)),
			commas(int(dtotal)), commas(int(isales)), commas(int(irfnd)),
			commas(int(itotal)), commas(int(gtotal)))
		dsalesT += dsales
		drfndT += drfnd
		isalesT += isales
		irfndT += irfnd
	}
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	fmt.Printf("%-10s %14s %14s %14s %14s %14s %14s %14s\n",
		"Total", commas(int(dsalesT)), commas(int(drfndT)),
		commas(int(dsalesT+drfndT)), commas(int(isalesT)),
		commas(int(irfndT)), commas(int(isalesT+irfndT)),
		commas(int(dsalesT+drfndT+isalesT+irfndT)))
	fmt.Printf("---------- -------------- -------------- -------------- -------------- -------------- -------------- --------------\n")
	fmt.Printf("\n")
}

// QuerySalesSQ show query results from database
func QuerySalesSQ(reportType string, fromDate string, toDate string, cfg Config) {
	fmt.Println(cfg.SqliteDb)
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch reportType {
	case "table":
		salesTabularSQ1(db, fromDate, toDate)
		salesTabularSQ2(db, fromDate, toDate)
	case "raw":
		salesRawSQ(db, fromDate, toDate)
	default:
		salesTabularSQ1(db, fromDate, toDate)
		salesTabularSQ2(db, fromDate, toDate)
	}
}

func commas(num int) string {
	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}
