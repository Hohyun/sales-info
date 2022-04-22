package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

func exportTabularSQ1(db *sql.DB, dstFile string, raw bool, vat bool, fromDate string, toDate string) {
	// Get data: Sale by date
	var sql string
	if raw {
		sql = fmt.Sprintf("select * from sales_yr_tax_raw where salesdate between '%s' and '%s'", fromDate, toDate)
	} else {
		if vat {
			sql = fmt.Sprintf("select * from sales_yr_tax_with_vat where salesdate between '%s' and '%s'", fromDate, toDate)
		} else {
			sql = fmt.Sprintf("select * from sales_yr_tax_without_vat where salesdate between '%s' and '%s'", fromDate, toDate)
		}
	}
	rows, err := db.Query(sql)
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

func exportTabularSQ2(db *sql.DB, dstFile string, raw bool, vat bool, fromDate string, toDate string) {
	var sql string
	if raw || vat {
		sql = fmt.Sprintf("select salesdate, dom_sales, dom_refund, dom_total, intl_sales, intl_refund, intl_total, g_total	from sales_refund where salesdate between '%s' and '%s'", fromDate, toDate)
	} else {
		sql = fmt.Sprintf("select salesdate, dom_sales_net, dom_refund_net, dom_total_net, intl_sales, intl_refund, intl_total, g_total_net from sales_refund where salesdate between '%s' and '%s'", fromDate, toDate)
	}
	rows, err := db.Query(sql)
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
func ExportCsvSQ(dstFile string, raw bool, vat bool, fromDate string, toDate string) {
	cfg := ParseConfig()
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	exportTabularSQ1(db, dstFile, raw, vat, fromDate, toDate)
	exportTabularSQ2(db, dstFile, raw, vat, fromDate, toDate)
}

