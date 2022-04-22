package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

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

// QuerySalesSQ show query results from database
func QuerySalesSQ(raw bool, vat bool, fromDate string, toDate string) {
	cfg := ParseConfig()
	db, err := sql.Open("sqlite3", cfg.SqliteDb)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	salesTabularSQ1(db, raw, vat, fromDate, toDate)
	salesTabularSQ2(db, raw, vat, fromDate, toDate)	
}

func salesTabularSQ1(db *sql.DB, raw bool, vat bool, fromDate string, toDate string) {
	var dsalesT, dyrT, dtaxT, dtotalT, isalesT, iyrT, itaxT, itotalT, gtotalT float64
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

	if raw {
		fmt.Printf("\n                                                                                                                                [ Raw Data ]\n")
	} else {
		if vat {
			fmt.Printf("\n                                                                                                                            [ VAT: included ]\n")
		} else {
			fmt.Printf("\n                                                                                                                            [ VAT: excluded ]\n")
		}
	}
	fmt.Printf("-------------------------------------------------------------------------------------------------------------------------------------------------\n")
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

func salesTabularSQ2(db *sql.DB, raw bool, vat bool, fromDate string, toDate string) {
	var dsalesT, drfndT, isalesT, irfndT float64
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

	if raw || vat {
		fmt.Printf("\n                                                                                                  [ VAT: included ]\n")
	} else {
		fmt.Printf("\n                                                                                                  [ VAT: excluded ]\n")
	}
	fmt.Printf("-------------------------------------------------------------------------------------------------------------------\n")
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

func commas(num int) string {
	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}
