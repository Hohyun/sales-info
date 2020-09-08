package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ConvertData(in_file string, out_file string) {

	f_in, err := os.Open(in_file)
	if err != nil {
		panic(err)
	}
	rdr := csv.NewReader(bufio.NewReader(f_in))
	records, err := rdr.ReadAll()
	if err != nil {
		panic(err)
	}

	results := [][]string{}
	//row := []string{}

	for _, r := range records {
		transformed := handle_row(r)
		// fmt.Println(transformed)
		results = append(results, transformed)
	}

	f_out, err := os.Create(out_file)
	wrt := csv.NewWriter(bufio.NewWriter(f_out))

	if err := wrt.WriteAll(results); err != nil {
		panic(err)
	}
	fmt.Println("File generated successfully --> " + out_file)
}

func handle_row(row []string) []string {
	fop := row[0]
	agency_type := row[2]
	sales_date := row[4]
	sales_type := row[5]
	ticket := row[7]
	itinerary := row[10]
	docs := row[13]
	ccy := row[14]
	amount := strings.ReplaceAll(row[15], ",", "")
	krw_amt := strings.ReplaceAll(row[16], ",", "")
	dom_intl := choose_di(fop, agency_type, sales_type, itinerary)
	sales_refund := choose_sr(krw_amt)
	return []string{fop, agency_type, sales_date, sales_type, ticket, itinerary, docs, ccy, amount, krw_amt, dom_intl, sales_refund}
}

func choose_di(fop string, agency_t string, sales_t string, itin string) string {
	// handle header line
	if fop == "FOP" {
		return "DomIntl"
	}

	switch sales_t {
	case "DOM":
		return "DOM"
	case "INTL":
		return "INTL"
	case "???":
		if strings.HasPrefix(agency_t, "BSP") || strings.HasPrefix(fop, "QN") || strings.Contains(itin, "ICN") {
			return "INTL"
		}
		if strings.Contains(itin, "GMP") {
			return "DOM"
		} else {
			return "INTL"
		}
	default:
		return "INTL"
	}
}

func choose_sr(krw_str string) string {
	// handle header line
	if krw_str == "KRW Amount" {
		return "KRW Amount"
	}

	i64, err := strconv.ParseInt(krw_str, 10, 64)
	if err != nil {
		panic(err)
	}
	if i64 >= 0 {
		return "Sales"
	} else {
		return "Refund"
	}
}
