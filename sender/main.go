package main

import (
	"flag"

	"os"
)

func main() {
	// get arguments
	host := *flag.String("host", "", "Mandatory - The host of the file relay")
	fn := *flag.String("file", "", "Mandatory - The name of the file you wish to transfer")

	flag.Parse()

	validateFlags(host, fn)

	// get file

	// checksum file

	// generate code

	// open connection with relay

	// begin transfer on accept, send checksum with meta data

	// wait for transfer to complete

	// no errors exit
}

func validateFlags(host string, fn string) {
	failed := false
	msg := "Error - Mandatory flag missing "

	if host == "" {
		println(msg + "'host'")
		failed = true
	}

	if fn == "" {
		println(msg + "'file'")
		failed = true
	}

	if failed {
		print("For help using this command tool please enter 'sender -h'")
		os.Exit(1)
	}
}
