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
	// fmt.Printf("in: %s, out: %s, src: %s, dst: %s, lvl: %d\n", flgIn, flgOut, flgSrc, flgDst, flgRpt)
	// fmt.Printf("args: %v", args)
	// fmt.Printf("%+v\n", c)

	if flgHelp || len(args) != 1 {
		DisplayUsage()
	}

	cmd := args[0]
	if cmd != "convert" && cmd != "import" && cmd != "export" && cmd != "query" {
		DisplayUsage()
	}

	backend := strings.ToLower(cfg.Database)
	switch cmd {
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
