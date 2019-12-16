package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/Samyoul/storj-file-sender/common"
	"github.com/Samyoul/storj-file-sender/sender/codegen"
)

func main() {
	// get arguments
	host := flag.String("host", "", "Mandatory - The host of the file relay")
	fn := flag.String("file", "", "Mandatory - The name of the file you wish to transfer")
	flag.Parse()

	err := validateFlags(host, fn)
	if err != nil {
		log.Fatalf("Error - Validating flags : %s", err)
	}

	// checksum file
	h, err := common.HashFile(*fn)
	if err != nil {
		log.Fatalf("Error - Checksumming file %s : %s", *fn, err)
	}

	// generate secret code
	// Use the int64 encoded checksum of the file as part of the random seed
	code := codegen.Make(int64(binary.BigEndian.Uint64(h.Sum(nil))))

	// display secret code
	println(code)

	// open connection with relay
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatalf("Error - making a connection : %s", err)
	}
	defer conn.Close()

	// Set write buffer size
	err = conn.(*net.TCPConn).SetWriteBuffer(common.BufferLimit)
	if err != nil {
		log.Fatalf("Error - setting write buffer : %s", err)
	}

	// Write and send header on connection, send checksum and filename with header
	hdr := common.MakeRequestHeaderSend(*fn, code, h.Sum(nil))
	conn.Write(hdr)

	// Open file hold ready to transfer
	f, err := os.Open(*fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// begin transfer on accept
	_, err = io.Copy(conn, f)
	if err != nil {
		log.Fatal(err)
	}

	// Add file terminator so the copy of connections knows to stop
	_, err = conn.Write(common.Terminator)
	if err != nil {
		log.Fatal(err)
	}

	// Read response
	buff := make([]byte, 2)
	_, err = conn.Read(buff)
	if err != nil {
		log.Fatal(err)
	}
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
