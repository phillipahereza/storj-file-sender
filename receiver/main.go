package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"github.com/Samyoul/storj-file-sender/common"
)

func main() {
	// get arguments
	args := os.Args
	err := validateArgs(args)
	if err != nil {
		log.Fatalf("error - validating arguments : %s", err)
	}

	// open connection with relay
	conn, err := net.Dial("tcp", args[1])
	if err != nil {
		log.Fatalf("error - making a connection : %s", err)
	}
	defer conn.Close()

	// make receive request to relay
	reqH := common.MakeRequestHeaderReceive(args[2])
	conn.Write(reqH)

	// get receive response header from relay with checksum and filename
	resH, err := common.GetResponseHeader(conn)
	if err != nil {
		log.Fatalf("error - reading response header : %s", err)
	}

	// start to receive data stream
	fn := args[3] + string(resH["Filename"])

	f, err := os.Create(fn)
	if err != nil {
		log.Fatalf("error - creating file : %s", err)
	}
	defer f.Close()

	// write data to file
	_, err = io.Copy(f, conn)
	if err != nil {
		log.Fatalf("error - creating file : %s", err)
	}

	// check file complete with checksum comparison.
	h, err := common.HashFile(fn)
	if err != nil {
		log.Fatalf("error - checksumming file %s : %s", fn, err)
	}

	if bytes.Compare(h.Sum(nil), resH["Checksum"]) != 0 {
		log.Fatalf("error - checksum does not match")
	}
}

func validateArgs(args []string) error {
	if len(args) != 4 {
		return errors.New(
			"invalid number of arguments.\n" +
				"expected : ./receiver <relay-host>:<relay-port> <secret-code> <output-directory>\n" +
				"example  : ./receiver localhost:9021 this-is-a-secret-code out/")
	}

	return nil
}
