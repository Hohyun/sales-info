package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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
)


// QuerySalesPG show query results from database
func QuerySalesPG(reportType string, fromDate string, toDate string) {
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
		salesTabularPG1(db, fromDate, toDate)
		salesTabularPG2(db, fromDate, toDate)
	case "raw":
		salesRawPG(db, fromDate, toDate)
	default:
		salesTabularPG1(db, fromDate, toDate)
		salesTabularPG2(db, fromDate, toDate)
	}
}

func salesRawPG(db *sql.DB, fromDate string, toDate string) {
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

func salesTabularPG1(db *sql.DB, fromDate string, toDate string) {
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

func salesTabularPG2(db *sql.DB, fromDate string, toDate string) {
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