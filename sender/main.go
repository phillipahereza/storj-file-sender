package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"github.com/Samyoul/storj-file-sender/common"
	"github.com/Samyoul/storj-file-sender/sender/codegen"
)

func main() {
	// get arguments
	args := os.Args

	err := validateArgs(args)
	if err != nil {
		log.Fatalf("Error - Validating arguments : %s", err)
	}

	// checksum file
	h, err := common.HashFile(args[2])
	if err != nil {
		log.Fatalf("Error - Checksumming file %s : %s", args[2], err)
	}

	// generate secret code
	// Use the int64 encoded checksum of the file as part of the random seed
	code := codegen.Make(int64(binary.BigEndian.Uint64(h.Sum(nil))))

	// display secret code
	println(code)

	// open connection with relay
	conn, err := net.Dial("tcp", args[1])
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
	hdr := common.MakeRequestHeaderSend(args[2], code, h.Sum(nil))
	conn.Write(hdr)

	// Open file hold ready to transfer
	f, err := os.Open(args[2])
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

func validateArgs(args []string) error {
	if len(args) != 3 {
		return errors.New(
			"mandatory arguments not present.\n" +
				"Expect the following arguments : `./send <relay-host>:<relay-port> <file-to-send>`\n" +
				"Example : `./send localhost:9021 corgis.mp4`")
	}

	return nil
}
