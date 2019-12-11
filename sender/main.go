package main

import (
	"flag"
	"os"
	"io"
	"crypto/sha256"
	"github.com/davecgh/go-spew/spew"
	"fmt"
	"encoding/hex"
)

func main() {
	// get arguments
	host := flag.String("host", "", "Mandatory - The host of the file relay")
	fn := flag.String("file", "", "Mandatory - The name of the file you wish to transfer")

	flag.Parse()
	validateFlags(host, fn)

	// checksum file
	cs, err := checksumFile(fn)
	if err != nil {
		fmt.Printf("Error - Checksumming file %s : %s\n", *fn, err)
		os.Exit(1)
	}

	// generate secret code

	// open connection with relay

	// display secret code

	// begin transfer on accept, send checksum with meta data

	// wait for transfer to complete

	// no errors exit
}

func validateFlags(host *string, fn *string) {
	failed := false
	msg := "Error - Mandatory flag missing "

	if *host == "" {
		println(msg + "'host'")
		failed = true
	}

	if *fn == "" {
		println(msg + "'file'")
		failed = true
	}

	if failed {
		print("For help using this command tool please enter 'sender -h'")
		os.Exit(1)
	}
}

func checksumFile(fn *string) (string, error) {
	h := sha256.New()

	f, err := os.Open(*fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
