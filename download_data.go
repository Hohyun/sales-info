package main

import (
	"fmt"
	"log"
	"os/exec"
)

func DownloadData(flgFrom string, flgTo string) {
	var vectisId, vectisPswd string
	if flgFrom == "" {
		fmt.Print("Input from date (yyyy-mm-dd): ")
		fmt.Scanln(&flgFrom)
	}
	if flgTo == "" {
		fmt.Print("Input to   date (yyyy-mm-dd): ")
		fmt.Scanln(&flgTo)
	}
	fmt.Print("Input Vectis ID:       ")
	fmt.Scanln(&vectisId)
	fmt.Print("Input Vectis Password: ")
	fmt.Scanln(&vectisPswd)

	fmt.Printf("from date: %s, to date: %s\n, id: %s, pasword: %s\n", flgFrom, flgTo, vectisId, vectisPswd)

	cmd := exec.Command("./sales_download.exe", flgFrom, flgTo, vectisId, vectisPswd)
	fmt.Printf("Downloading sales data for %s ~ %s ....\n", flgFrom, flgTo)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Download failed with '%s'\n", err)
		return
	}
}
