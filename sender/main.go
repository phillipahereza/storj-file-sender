package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"errors"
)

func main() {
	// get arguments
	host := flag.String("host", "", "Mandatory - The host of the file relay")
	fn := flag.String("file", "", "Mandatory - The name of the file you wish to transfer")

	flag.Parse()
	err := validateFlags(host, fn)
	if err != nil {
		fmt.Printf("Error - Validating flags : %s", err)
		os.Exit(1)
	}

	// checksum file
	_, err = checksumFile(fn)
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

func validateFlags(host *string, fn *string) error {
	failed := false
	msg := "Mandatory flag(s) missing : "

	if *host == "" {
		msg += "'host', "
		failed = true
	}

	if *fn == "" {
		msg += "'file', "
		failed = true
	}

	if failed {
		msg += "\nFor help using this command tool please enter 'sender -h'"
		return errors.New(msg)
	}

	return nil
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
