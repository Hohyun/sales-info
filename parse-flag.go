package main

import (
	"encoding/json"
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
	flgLvl  int
)

// Config is data struct for parsing config.json
type Config struct {
	RootDir string `json:"root_dir"`
	DataDir string `json:"data_dir"`
}

// ParseCmdLineFlags is function parsing command line argument options
func ParseCmdLineFlags() {
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgIn, "in", path.Join(getDataDir(), "vectis_sales.csv"), "convert: input file")
	flag.StringVar(&flgOut, "out", path.Join(getDataDir(), "sales_for_db.csv"), "convert: output file")
	flag.StringVar(&flgSrc, "src", path.Join(getDataDir(), "sales_for_db.csv"), "import: source file")
	flag.StringVar(&flgDst, "dst", path.Join(getDataDir(), "sales_info.xlsx"), "export: output file")
	flag.IntVar(&flgLvl, "lvl", 0, "level (1: Grand total, 2:D/I, 3:D/I+S/R, 4:Date+D/I+S/R)")
	flag.Parse()
}

// DisplayUsage shows how to use program.
func DisplayUsage() {
	fmt.Println("Usage: sales-info [-lvl 1|2|3|4 ] query  | [-in  filename -out filename] convert |")
	fmt.Println("                  [-src filename] import | [-dst filename] export")
	flag.PrintDefaults()
	os.Exit(0)
}

func getDataDir() string {
	f, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var c Config
	dec := json.NewDecoder(f)
	dec.Decode(&c)
	// fmt.Printf("%+v\n", c)
	return c.DataDir
}
