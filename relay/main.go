package main

import (
	"flag"
	"os"
)

func main() {

	// get argument
	port := flag.Int("port", 0, "Mandatory - The port to listen on")

	flag.Parse()
	validateFlags(port)

	// open TCP server

	// wait for connections

		// serve connections with relay() ... will probably need a handler for each the send and receive requests
}

func relay() {
	// get incoming send relay request

	// wait for relay receive request

	// begin data stream

	// on completion close connection
}

func validateFlags(port *int) {
	failed := false
	msg := "Error - Mandatory flag missing "

	if *port == 0 {
		println(msg + "'port'")
		failed = true
	}

	if failed {
		print("For help using this tool please enter 'relay -h'")
		os.Exit(1)
	}
}