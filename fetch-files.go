package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// FetchFiles : fetch invoice and payment files from ftp server
func FetchFiles(from string, to string) {
	fmt.Println("Fetch files from ftp server...")
	addr := "10.23.34.4:22"
	config := &ssh.ClientConfig{
		User: "ibsSale2",
		Auth: []ssh.AuthMethod{
			ssh.Password(""),
		},
		HostKeyCallback: trustedHostKeyCallback("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEAwBwnwZeKacxrPQEo8UGRiIlS5UQR0tYYtSfD6tjFhAaBhUA5BF39f4XaKr2hhU7K3ZVBJP+1pldwIDnCCuNksH5EBiRKkHB46CKVZWlE/GbH0jgWkZARzXsNNGx+jAtaPU7LkljQnPj8y0/ucAruQcFhOSldaykny5a4ppLTgI0="), // <- server-key goes here

		Timeout: 30 * time.Second,
		Config: ssh.Config{
			// ciplers: aes256-cbc
			Ciphers:      []string{"aes128-cbc", "3des-cbc", "blowfish-cbc", "aes192-cbc", "aes256-cbc", "rijndael128-cbc", "rijndael192-cbc", "rijndael256-cbc", "rijndael-cbc@lysator.liu.se"},
			KeyExchanges: []string{"diffie-hellman-group1-sha1"},
		},
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("failed to connect sftp server: %v", err)
	}
	defer conn.Close()

	// Create new SFTP client
	sc, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatalf("Unable to start SFTP subsystem: %v\n", err)
	}
	defer sc.Close()
	downloadFiles(sc, from, to)
}

func downloadFiles(sc *sftp.Client, from string, to string) (err error) {
	remoteDir := "/SalesData"
	log.Printf("Remote Directory [%s] ...", remoteDir)

	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		log.Fatalf("Unable to list remote dir: %v", err)
		return
	}

	// default date: yesterday
	if from == "" || to == "" {
		dt := time.Now().AddDate(0, 0, -1)
		year, month, day := dt.Date()
		from = fmt.Sprintf("%d-%02d-%02d", year, month, day)
		to = fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}

	for _, f := range files {
		yymmdd1 := strings.ReplaceAll(from[2:], "-", "")
		yymmdd2 := strings.ReplaceAll(to[2:], "-", "")
		fopPrefix := "SALE_FOP_TKT_" + yymmdd1 + "-" + yymmdd2
		taxPrefix := "SALE_TAX_TKT_" + yymmdd1 + "-" + yymmdd2
		filename := f.Name()

		if strings.HasPrefix(filename, fopPrefix) || strings.HasPrefix(filename, taxPrefix) {
			downloadAction(sc, remoteDir, filename)
		} 
	}

	return
}

// Download file from sftp server
func downloadAction(sc *sftp.Client, remoteDir string, remoteFile string) (err error) {
	localFile := ""

	if strings.HasPrefix(remoteFile, "SALE_FOP") {
		localFile = "./data/" + "VectisReport_sales.csv"
	} else if strings.HasPrefix(remoteFile, "SALE_TAX"){
		localFile = "./data/" + "VectisReport_taxyr.csv"
	}

	log.Printf("Downloading: %s/%s ...", remoteDir, remoteFile)
	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := sc.OpenFile(remoteDir+"/"+remoteFile, (os.O_RDONLY))
	if err != nil {
		log.Fatalf("Unable to open remote file: %v", err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFile)
	if err != nil {
		log.Fatalf("Unable to open local file: %v", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatalf("Unable to download remote file: %v", err)
	}

	log.Printf("Done: ----> %s", localFile)
	return
}

// create human-readable SSH-key strings
func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}

func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {
	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			log.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}
		return nil
	}
}
