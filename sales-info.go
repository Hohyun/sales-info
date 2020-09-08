package main

import (
	"flag"
	"fmt"
	//"log"
	//"os"
)

func main() {
	ParseCmdLineFlags()
	args := flag.Args()

	fmt.Printf("in: %s, out: %s, src: %s, dst: %s, lvl: %d\n", flgIn, flgOut, flgSrc, flgDst, flgLvl)
	fmt.Println(args)

	if flgHelp || len(args) !=1 {
		DisplayUsage()
	}

	cmd := args[0]
	if cmd != "convert" && cmd != "import" && cmd != "export" && cmd != "query" {
		DisplayUsage()
	}

	switch cmd {
	case "convert": ConvertData(flgIn, flgOut)
	case "import": fmt.Println("Import undefined")
	case "export": fmt.Println("Export undefined")
	case "query": QuerySales()
	default: fmt.Println("What happen? I don't know.")
		
	}
}
