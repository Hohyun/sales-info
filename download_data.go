package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	//"os/exec"
	"path"
)

// DownloadData : excute vectis download external program.
func DownloadData(flgFrom string, flgTo string, flgID string, flgPswd string, cfg Config) {
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

	cerr := copyFileContents(cfg.Data.VectisFile, path.Join(cfg.Data.DirName, cfg.Data.SourceFile))
	if cerr != nil {
		log.Fatalf("Failed to copy vectis report to %s : '%s'\n", cfg.Data.DirName, cerr)
		return
	}
	fmt.Printf("Download file was copied to %s\n", path.Join(cfg.Data.DirName, cfg.Data.SourceFile))
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
