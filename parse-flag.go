package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	flgHelp bool
	flgFrom string
	flgTo   string
	flgIn   string
	flgOut  string
	flgSrc  string
	flgDst  string
	flgRpt  string
)

// ParseCmdLineFlags is function parsing command line argument options
func ParseCmdLineFlags(cfg Config) {
	d := cfg.Data.DirName
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgFrom, "from", "", ":download: from date (yyyy-mm-dd)")
	flag.StringVar(&flgTo, "to", "", ":download: to date (yyyy-mm-dd) ")
	flag.StringVar(&flgIn, "in", path.Join(d, cfg.Data.SourceFile), ":convert : input file ")
	flag.StringVar(&flgOut, "out", path.Join(d, cfg.Data.ImportFile), ":convert : output file")
	flag.StringVar(&flgSrc, "src", path.Join(d, cfg.Data.ImportFile), ":import  : source file")
	flag.StringVar(&flgDst, "dst", path.Join(d, cfg.Data.ExportFile), ":export  : output file")
	flag.StringVar(&flgRpt, "rpt", "table", ":query  : report type")
	flag.Parse()
}

// DisplayUsage shows how to use program.
func DisplayUsage() {
	fmt.Println("Usage: sales-info -from yyyy-mm-dd -to yyyy-mm-dd download |")
	fmt.Println("       [-in filename -out filename] convert | [-src filename] import |")
	fmt.Println("       [-rpt table|raw] query | [-dst filename] export | all")
	flag.PrintDefaults()
	os.Exit(0)
}
