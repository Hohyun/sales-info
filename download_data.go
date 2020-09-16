package main

import (
	"fmt"
	"log"
	"os/exec"
)

// DownloadData : excute vectis download external program.
func DownloadData(flgFrom string, flgTo string, flgID string, flgPswd string) {
	if flgFrom == "" {
		fmt.Print("Input from date (yyyy-mm-dd): ")
		fmt.Scanln(&flgFrom)
	}
	if flgTo == "" {
		fmt.Print("Input to   date (yyyy-mm-dd): ")
		fmt.Scanln(&flgTo)
	}
	if flgID == "" {
		fmt.Print("Input Vectis ID:       ")
		fmt.Scanln(&flgID)
	}
	if flgPswd == "" {
		fmt.Print("Input Vectis Password: ")
		fmt.Scanln(&flgPswd)
	}
	// fmt.Printf("from date: %s, to date: %s\n, id: %s, pasword: %s\n", flgFrom, flgTo, flgID, flgPswd)
	cmd := exec.Command("./sales_download.exe", flgFrom, flgTo, flgID, flgPswd)
	fmt.Printf("Downloading sales data for %s ~ %s ....\n", flgFrom, flgTo)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Download failed with '%s'\n", err)
		return
	}
}
