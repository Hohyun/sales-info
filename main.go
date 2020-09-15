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

	if flgHelp || len(args) != 1 {
		DisplayUsage()
	}

	backend := strings.ToLower(cfg.Database)
	action := args[0]
	switch action {
	case "download":
		DownloadData(flgFrom, flgTo)
	case "all":
		ConvertData(flgIn, flgOut)
		if backend == "postgresql" {
			ImportCsvPG(flgSrc)
			QuerySalesPG(flgRpt)
			ExportCsvPG(flgDst)
		} else {
			ImportCsvSQ(flgSrc, cfg)
			QuerySalesSQ(flgRpt, cfg)
			ExportCsvSQ(flgDst, cfg)
		}
	case "convert":
		ConvertData(flgIn, flgOut)
	case "import":
		if backend == "postgresql" {
			ImportCsvPG(flgSrc)
		} else {
			ImportCsvSQ(flgSrc, cfg)
		}
	case "export":
		if backend == "postgresql" {
			ExportCsvPG(flgDst)
		} else {
			ExportCsvSQ(flgDst, cfg)
		}
	case "query":
		if backend == "postgresql" {
			QuerySalesPG(flgRpt)
		} else {
			QuerySalesSQ(flgRpt, cfg)
		}
	default:
		DisplayUsage()
	}
	fmt.Println("Bye!")
}
