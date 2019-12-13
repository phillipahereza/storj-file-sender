package main

import (
	"flag"
	"os"
)

func main() {
	// get arguments
	host := flag.String("host", "", "Mandatory - The host of the file relay")
	code := flag.String("code", "", "Mandatory - The secret code of the file you wish to receive")
	dir := flag.String("out", "", "Mandatory - The name of the directory name you wish to receive the file to")
	flag.Parse()

	validateFlags(host, code, dir)

	// make receive request to relay

	// start to receive data stream and checksum in meta data

	// write data to file from arguments

	// check file complete with checksum comparison.

	// no errors exit
}

func validateFlags(host *string, code *string, dir *string) {
	failed := false
	msg := "Error - Mandatory flag missing "

	if *host == "" {
		println(msg + "'host'")
		failed = true
	}

	if *code == "" {
		println(msg + "'code'")
		failed = true
	}

	if *dir == "" {
		println(msg + "'out'")
		failed = true
	}

	if failed {
		print("For help using this command tool please enter 'receiver -h'")
		os.Exit(1)
	}
}
