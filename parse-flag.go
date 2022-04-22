package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	flgHelp  bool
	flgGubun string
	flgFrom  string
	flgTo    string
	flgID    string
	flgPswd  string
	flgIn    string
	flgOut   string
	flgSrc   string
	flgDst   string
	flgVat   bool
	flgRaw   bool
)

// ParseCmdLineFlags is function parsing command line argument options
func ParseCmdLineFlags(cfg Config) {
	d := cfg.Data.DirName
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgGubun, "gubun", "", "sales or taxyr")
	flag.StringVar(&flgFrom, "from", "", "from date (yyyy-mm-dd, default: yesterday)")
	flag.StringVar(&flgTo, "to", "", "to date (yyyy-mm-dd, default: yesterday) ")
	flag.StringVar(&flgID, "id", "", "vectis id")
	flag.StringVar(&flgPswd, "pswd", "", "vectis password")
	flag.StringVar(&flgIn, "in", path.Join(d, cfg.Data.SourceFile), "input file ")
	flag.StringVar(&flgOut, "out", path.Join(d, cfg.Data.ImportFile), "output file")
	flag.StringVar(&flgSrc, "src", path.Join(d, cfg.Data.ImportFile), "source file")
	flag.StringVar(&flgDst, "dst", path.Join(d, cfg.Data.ExportFile), "output file")
	flag.BoolVar(&flgVat, "vat", false, "include VAT in results")
	flag.BoolVar(&flgRaw, "raw", false, "report for raw data")
	flag.Parse()
}

// DisplayUsage shows how to use program.
func DisplayUsage() {
	fmt.Println("Usage: sales-info")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd                                            ")
	fmt.Println("         -gubun sales|taxyr -id ******* -pswd ********                  ] download |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd                                ] fetch    |")
	fmt.Println("       [ -gubun sales|taxyr -in filename -out filename                  ] convert  |")
	fmt.Println("       [ -gubun sales|taxyr -src filename                               ] import   |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -raw -vat                      ] query    |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -raw -vat -dst filename        ] export   |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -raw -vat -dst filename        ] all       ")
	fmt.Println(" * all : fetch -> convert -> import -> query -> export                              ")
	flag.PrintDefaults()
	os.Exit(0)
}
