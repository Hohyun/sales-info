package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	//"os/exec"
	"path"
)

// DownloadData : excute vectis download external program.
func DownloadData(flgGubun string, flgFrom string, flgTo string, flgID string, flgPswd string, cfg Config) {
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

	if flgGubun == "" {
		DownloadDataSub("sales", flgFrom, flgTo, flgID, flgPswd, cfg)
		DownloadDataSub("taxyr", flgFrom, flgTo, flgID, flgPswd, cfg)
	} else {
		DownloadDataSub(flgGubun, flgFrom, flgTo, flgID, flgPswd, cfg)
	}
}

// DownloadData : excute vectis download external program.
func DownloadDataSub(gubun string, flgFrom string, flgTo string, flgID string, flgPswd string, cfg Config) {
	var programName, dstFile string
	if gubun == "sales" {
		programName = "sales_download.exe"
		dstFile = strings.Replace(cfg.Data.SourceFile, ".", "_sales.", 1)
	} else if gubun == "taxyr" {
		programName = "taxyr_download.exe"
		dstFile = strings.Replace(cfg.Data.SourceFile, ".", "_taxyr.", 1)
	} 

	fmt.Printf("Downloading %s data for %s ~ %s ....\n", gubun, flgFrom, flgTo)
	cmd := exec.Command(programName, flgFrom, flgTo, flgID, flgPswd)

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Wait() failed with '%s'\n", err)
	}

	err = copyFileContents(cfg.Data.VectisFile, path.Join(cfg.Data.DirName, dstFile))
	if err != nil {
		log.Fatalf("Failed to copy vectis report to %s : '%s'\n", cfg.Data.DirName, err)
		return
	}
	fmt.Printf("Download file was copied to %s\n", path.Join(cfg.Data.DirName, dstFile))
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
