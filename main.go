package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	cfg := ParseConfig()
	ParseCmdLineFlags(cfg)
	args := flag.Args()
	// fmt.Printf("from: %s, to: %s, in: %s, out: %s, src: %s, dst: %s, lvl: %d\n", flgFrom, flgTo, flgIn, flgOut, flgSrc, flgDst, flgRpt)
	// fmt.Printf("args: %v", args)
	// fmt.Printf("%+v\n", c)

	if flgHelp || len(args) != 1 { // there should be only one action command.
		DisplayUsage()
	}
	if flgFrom == "" {
		flgFrom = getDefautFromDate()
	}
	if flgTo == "" {
		flgTo = getDefautToDate()
	}
	action := args[0]

	switch action {
	case "download":
		DownloadData(flgGubun, flgFrom, flgTo, flgID, flgPswd, cfg)
	case "fetch":
		FetchFiles(flgFrom, flgTo)
	case "all":
		FetchFiles(flgFrom, flgTo)
		ConvertData("sales", strings.Replace(flgIn, ".", "_sales.", 1), strings.Replace(flgOut, ".", "_sales.", 1))
		ConvertData("taxyr", strings.Replace(flgIn, ".", "_taxyr.", 1), strings.Replace(flgOut, ".", "_taxyr.", 1))
		if flgDB == "pg" {
			if flgGubun == "sales" {
				ImportCsvPgSales(strings.Replace(flgSrc, ".", "_sales.", 1))
			} else if flgGubun == "taxyr" {
				ImportCsvPgTaxYr(strings.Replace(flgSrc, ".", "_taxyr.", 1))
			} else if flgGubun == "" {
				ImportCsvPgSales(strings.Replace(flgSrc, ".", "_sales.", 1))
				ImportCsvPgTaxYr(strings.Replace(flgSrc, ".", "_taxyr.", 1))
			}
			QuerySalesPG(true, false, flgFrom, flgTo)
			ExportCsvPG(flgDst, true, false, flgFrom, flgTo)
		} else {
			if flgGubun == "sales" {
				ImportCsvSqSales(cfg)
			} else if flgGubun == "taxyr" {
				ImportCsvSqTaxYr(cfg)
			} else if flgGubun == "" {
				ImportCsvSqSales(cfg)
				ImportCsvSqTaxYr(cfg)
			}
			QuerySalesSQ(true, false, flgFrom, flgTo)
			ExportCsvSQ(flgDst, true, false, flgFrom, flgTo)
		}
	case "convert":
		if flgGubun == "sales" {
			ConvertData("sales", strings.Replace(flgIn, ".", "_sales.", 1), strings.Replace(flgOut, ".", "_sales.", 1))
		} else if flgGubun == "taxyr" {
			ConvertData("taxyr", strings.Replace(flgIn, ".", "_taxyr.", 1), strings.Replace(flgOut, ".", "_taxyr.", 1))
		} else if flgGubun == "" {
			ConvertData("sales", strings.Replace(flgIn, ".", "_sales.", 1), strings.Replace(flgOut, ".", "_sales.", 1))
			ConvertData("taxyr", strings.Replace(flgIn, ".", "_taxyr.", 1), strings.Replace(flgOut, ".", "_taxyr.", 1))
		}

	case "import":
		if flgDB == "pg" {
			if flgGubun == "sales" {
				ImportCsvPgSales(strings.Replace(flgSrc, ".", "_sales.", 1))
			} else if flgGubun == "taxyr" {
				ImportCsvPgTaxYr(strings.Replace(flgSrc, ".", "_taxyr.", 1))
			} else if flgGubun == "" {
				ImportCsvPgSales(strings.Replace(flgSrc, ".", "_sales.", 1))
				ImportCsvPgTaxYr(strings.Replace(flgSrc, ".", "_taxyr.", 1))
			}
		} else {
			if flgGubun == "sales" {
				ImportCsvSqSales(cfg)
			} else if flgGubun == "taxyr" {
				ImportCsvSqTaxYr(cfg)
			} else if flgGubun == "" {
				ImportCsvSqSales(cfg)
				ImportCsvSqTaxYr(cfg)
			}
		}
	case "export":
		if flgDB == "pg" {
			ExportCsvPG(flgDst, flgRaw, flgVat, flgFrom, flgTo)
		} else {
			ExportCsvSQ(flgDst, flgRaw, flgVat, flgFrom, flgTo)
		}
	case "query":
		if flgDB == "pg" {
			QuerySalesPG(flgRaw, flgVat, flgFrom, flgTo)
		} else {
			QuerySalesSQ(flgRaw, flgVat, flgFrom, flgTo)
		}
	default:
		DisplayUsage()
	}
	fmt.Println("Bye!")
}
