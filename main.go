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

	backend := strings.ToLower(cfg.Database)
	action := args[0]
	if (action == "download" || action == "query" || action == "export") &&
		(flgFrom == "" || flgTo == "") {
		fmt.Printf("From date, To date should be supplied with %s\n", action)
		DisplayUsage()
	}

	switch action {
	case "download":
		DownloadData(flgGubun, flgFrom, flgTo, flgID, flgPswd, cfg)
	case "all":
		ConvertData("sales", strings.Replace(flgIn, ".", "_sales.", 1), strings.Replace(flgOut, ".", "_sales.", 1))
		ConvertData("taxyr", strings.Replace(flgIn, ".", "_taxyr.", 1), strings.Replace(flgOut, ".", "_taxyr.", 1))
		if backend == "postgresql" {
			ImportCsvPG(flgSrc)
			QuerySalesPG(flgRpt, flgFrom, flgTo)
			ExportCsvPG(flgRpt, flgDst, flgFrom, flgTo)
		} else {
			if flgGubun == "sales" {
				ImportCsvSqSales(cfg)
			} else if flgGubun == "taxyr" {
				ImportCsvSqTaxYr(cfg)
			} else if flgGubun == "" {
				ImportCsvSqSales(cfg)
				ImportCsvSqTaxYr(cfg)
			}	
			QuerySalesSQ(flgRpt, flgFrom, flgTo, cfg)
			ExportCsvSQ(flgRpt, flgDst, flgFrom, flgTo, cfg)
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
		if backend == "postgresql" {
			ImportCsvPG(flgSrc)
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
		if backend == "postgresql" {
			ExportCsvPG(flgRpt, flgDst, flgFrom, flgTo)
		} else {
			ExportCsvSQ(flgRpt, flgDst, flgFrom, flgTo, cfg)
		}
	case "query":
		if backend == "postgresql" {
			QuerySalesPG(flgRpt, flgFrom, flgTo)
		} else {
			QuerySalesSQ(flgRpt, flgFrom, flgTo, cfg)
		}
	default:
		DisplayUsage()
	}
	fmt.Println("Bye!")
}
