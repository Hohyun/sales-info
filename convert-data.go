package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ConvertData convert raw vectics file into csv for db import
func ConvertData(inFile string, outFile string) {

	fIn, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	rdr := csv.NewReader(bufio.NewReader(fIn))
	records, err := rdr.ReadAll()
	if err != nil {
		panic(err)
	}

	results := [][]string{}
	//row := []string{}

	for _, r := range records {
		transformed := handleRow(r)
		// fmt.Println(transformed)
		results = append(results, transformed)
	}

	fOut, err := os.Create(outFile)
	wrt := csv.NewWriter(bufio.NewWriter(fOut))

	if err := wrt.WriteAll(results); err != nil {
		panic(err)
	}
	fmt.Println("File generated successfully --> " + outFile)
}

func handleRow(row []string) []string {
	fop := row[0]
	fopdesc := row[1]
	agencyType := row[2]
	salesDate := row[4]
	salesType := row[5]
	ticket := row[7]
	itinerary := row[10]
	docs := row[13]
	ccy := row[14]
	amount := strings.ReplaceAll(row[15], ",", "")
	krwAmt := strings.ReplaceAll(row[16], ",", "")
	domIntl := chooseDi(fop, agencyType, salesType, itinerary)
	salesRefund := chooseSr(krwAmt)
	return []string{fop, fopdesc, agencyType, salesDate, salesType, ticket, itinerary, docs, ccy, amount, krwAmt, domIntl, salesRefund}
}

func chooseDi(fop string, agencyT string, salesT string, itin string) string {
	// handle header line
	if fop == "FOP" {
		return "DomIntl"
	}

	switch salesT {
	case "DOM":
		return "DOM"
	case "INTL":
		return "INTL"
	case "???":
		if strings.HasPrefix(agencyT, "BSP") || strings.HasPrefix(fop, "QN") || strings.Contains(itin, "ICN") {
			return "INTL"
		}
		if strings.Contains(itin, "GMP") {
			return "DOM"
		}
		return "INTL"
	default:
		return "INTL"
	}
}

func chooseSr(krwStr string) string {
	// handle header line
	if krwStr == "KRW Amount" {
		return "KRW Amount"
	}

	i64, err := strconv.ParseInt(krwStr, 10, 64)
	if err != nil {
		panic(err)
	}
	if i64 >= 0 {
		return "Sales"
	}
	return "Refund"
}
