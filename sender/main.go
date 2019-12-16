package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"hash"
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
	h, err := hashFile(fn)
	if err != nil {
		log.Fatalf("Error - Checksumming file %s : %s\n", *fn, err)
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

	// Write and send header on connection send checksum with header
	hdr := common.MakeRequestHeaderSend(code, h.Sum(nil))
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

func hashFile(fn *string) (hash.Hash, error) {
	h := sha256.New()

	f, err := os.Open(*fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h, nil
}
