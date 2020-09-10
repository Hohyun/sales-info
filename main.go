package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	ParseCmdLineFlags()
	args := flag.Args()
	c := ParseConfig()
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

	backend := strings.ToLower(c.Database)
	switch cmd {
	case "convert":
		ConvertData(flgIn, flgOut)
	case "import":
		if backend == "postgresql" {
			ImportCsvPG(flgSrc)
		} else {

			// ImportCsvSQ(flgSrc)
		}
	case "export":
		if backend == "postgresql" {
			ExportCsvPG(flgDst)
		} else {
			// ExportCsvSQ(flgDst)
		}
	case "query":
		if backend == "postgresql" {
			QuerySalesPG(flgRpt)
		} else {
			// QuerySalesSQ(flgRpt)
		}
	default:
		DisplayUsage()
	}
	fmt.Println("Bye!")
}
