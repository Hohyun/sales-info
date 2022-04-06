package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func exportRawPG(db *sql.DB, dstFile string, fromDate string, toDate string) {
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

func exportTabularPG1(db *sql.DB, dstFile string, fromDate string, toDate string) {
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

func exportTabularPG2(db *sql.DB, dstFile string, fromDate string, toDate string) {
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

// ExportCsvPG : export sales data with different format.
func ExportCsvPG(reportType string, dstFile string, fromDate string, toDate string) {
	// Open database
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	switch reportType {
	case "tabular":
		exportTabularPG1(db, dstFile, fromDate, toDate)
		exportTabularPG2(db, dstFile, fromDate, toDate)
	case "raw":
		exportRawPG(db, dstFile, fromDate, toDate)
	default:
		exportTabularPG1(db, dstFile, fromDate, toDate)
		exportTabularPG2(db, dstFile, fromDate, toDate)
	}
}
