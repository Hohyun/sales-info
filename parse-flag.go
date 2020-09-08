package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flgHelp bool
	flgIn string
	flgOut string 
	flgSrc string
	flgDst string
	flgLvl int
)

func ParseCmdLineFlags() {
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgIn, "in", "./data/vectis_sales.csv", "input file path")
	flag.StringVar(&flgOut, "out", "./data/sales.csv", "output file path")
	flag.StringVar(&flgSrc, "src", "./data/sales.csv", "source file path")
	flag.StringVar(&flgDst, "dst", "./data/sales.xlsx", "destination file path")
	flag.IntVar(&flgLvl, "lvl", 0, "level (0: All, 1:D/I, 2:D/I+S/R, 3:Date+D/I+S/R)")
	flag.Parse()
}

func DisplayUsage() {
	fmt.Println("Examples:")
	fmt.Println("sales-info [-in filename -out filename] convert")
	fmt.Println("sales-info [-src filename] import")
	fmt.Println("sales-info [-dst filename] export")
	fmt.Printf("sales-info [-lvl |0|1|2|3|] query\n\n")
	flag.Usage()
	os.Exit(0)
}