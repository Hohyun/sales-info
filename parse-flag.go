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

type Config struct {
	RootDir string `json:"root_dir"`
	DataDir string `json:"data_dir"`
}

func ParseCmdLineFlags() {
	flag.BoolVar(&flgHelp, "help", false, "show help")
	flag.StringVar(&flgIn, "in", path.Join(getDataDir(), "vectis_sales.csv"), "convert: input file")
	flag.StringVar(&flgOut, "out", path.Join(getDataDir(), "sales_for_db.csv"), "convert: output file")
	flag.StringVar(&flgSrc, "src", path.Join(getDataDir(), "sales_for_db.csv"), "import: source file")
	flag.StringVar(&flgDst, "dst", path.Join(getDataDir(), "sales_info.xlsx"), "export: output file")
	flag.IntVar(&flgLvl, "lvl", 4, "level (0:Grand total, 1:D/I, 2:D/I+S/R, 3:Date+D/I+S/R, 4:All)")
	flag.Parse()
}

func DisplayUsage() {
	fmt.Println("Usage:")
	fmt.Println("sales-info [-in filename -out filename] convert")
	fmt.Println("sales-info [-src filename] import")
	fmt.Println("sales-info [-dst filename] export")
	fmt.Printf("sales-info [-lvl |0|1|2|3|] query\n\n")
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
