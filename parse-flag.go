package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	flgHelp bool
	flgIn   string
	flgOut  string
	flgSrc  string
	flgDst  string
	flgRpt  string
)

// ParseCmdLineFlags is function parsing command line argument options
func ParseCmdLineFlags() {
	d := ParseConfig().DataDir
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgIn, "in", path.Join(d, "vectis_sales.csv"), ":convert: input file ")
	flag.StringVar(&flgOut, "out", path.Join(d, "sales_for_db.csv"), ":convert: output file")
	flag.StringVar(&flgSrc, "src", path.Join(d, "sales_for_db.csv"), ":import : source file")
	flag.StringVar(&flgDst, "dst", path.Join(d, "sales_info.csv"), ":export : output file")
	flag.StringVar(&flgRpt, "rpt", "table", ":query  : report type")
	flag.Parse()
}

// DisplayUsage shows how to use program.
func DisplayUsage() {
	fmt.Println("Usage: sales-info [-rpt table|raw] query  | [-in  filename -out filename] convert |")
	fmt.Println("                  [-src filename ] import | [-dst filename] export")
	flag.PrintDefaults()
	os.Exit(0)
}
