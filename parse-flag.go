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
	flgDate  string
	flgFrom  string
	flgTo    string
	flgID    string
	flgPswd  string
	flgIn    string
	flgOut   string
	flgSrc   string
	flgDst   string
	flgRpt   string
	flgVat   bool
)

// ParseCmdLineFlags is function parsing command line argument options
func ParseCmdLineFlags(cfg Config) {
	d := cfg.Data.DirName
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgGubun, "gubun", "", "sales or taxyr")
	flag.StringVar(&flgDate, "date", "", "sales date (yyyy-mm-dd, default: yesterday)")
	flag.StringVar(&flgFrom, "from", "", "from date (yyyy-mm-dd, default: yesterday)")
	flag.StringVar(&flgTo, "to", "", "to date (yyyy-mm-dd, default: yesterday) ")
	flag.StringVar(&flgID, "id", "", "vectis id")
	flag.StringVar(&flgPswd, "pswd", "", "vectis password")
	flag.StringVar(&flgIn, "in", path.Join(d, cfg.Data.SourceFile), "input file ")
	flag.StringVar(&flgOut, "out", path.Join(d, cfg.Data.ImportFile), "output file")
	flag.StringVar(&flgSrc, "src", path.Join(d, cfg.Data.ImportFile), "source file")
	flag.StringVar(&flgDst, "dst", path.Join(d, cfg.Data.ExportFile), "output file")
	flag.StringVar(&flgRpt, "rpt", "tabular", "report type")
	flag.BoolVar(&flgVat, "vat", false, "include VAT in results")
	flag.Parse()
}

// DisplayUsage shows how to use program.
func DisplayUsage() {
	fmt.Println("Usage: sales-info")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd                                            ")
	fmt.Println("         -gubun sales|taxyr -id ******* -pswd ********                  ] download |")
	fmt.Println("       [ -date yyyy-mm-dd                                               ] fetch    |")
	fmt.Println("       [ -gubun sales|taxyr -in filename -out filename                  ] convert  |")
	fmt.Println("       [ -gubun sales|taxyr -src filename                               ] import   |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -rpt tabular|raw -vat          ] query    |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -rtp tabular|raw -dst filename ] export   |")
	fmt.Println("       [ -from yyyy-mm-dd -to yyyy-mm-dd -date yyyy-mm-dd               ] all       ")
	fmt.Println(" * all : fetch -> convert -> import -> query -> export                              ")
	flag.PrintDefaults()
	os.Exit(0)
}
