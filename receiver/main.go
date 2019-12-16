package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/Samyoul/storj-file-sender/common"
)

func main() {
	// get arguments
	host := flag.String("host", "", "Mandatory - The host of the file relay")
	code := flag.String("code", "", "Mandatory - The secret code of the file you wish to receive")
	dir := flag.String("out", "./", "Mandatory - The name of the directory name you wish to receive the file to")
	flag.Parse()

	validateFlags(host, code, dir)

	// open connection with relay
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		log.Fatalf("Error - making a connection : %s", err)
	}
	defer conn.Close()

	// make receive request to relay
	reqH := common.MakeRequestHeaderReceive(*code)
	conn.Write(reqH)

	// get receive response header from relay with checksum and filename
	resH, err := common.GetResponseHeader(conn)
	if err != nil {
		log.Fatalf("Error - reading response header : %s", err)
	}

	// start to receive data stream
	fn := *dir + string(resH["Filename"])

	f, err := os.Create(fn)
	if err != nil {
		log.Fatalf("Error - creating file : %s", err)
	}
	defer f.Close()

	// write data to file
	_, err = io.Copy(f, conn)
	if err != nil {
		log.Fatalf("Error - creating file : %s", err)
	}

	// check file complete with checksum comparison.
	h, err := common.HashFile(fn)
	if err != nil {
		log.Fatalf("Error - Checksumming file %s : %s", fn, err)
	}

	if bytes.Compare(h.Sum(nil), resH["Checksum"]) != 0 {
		log.Fatalf("Error - Checksum does not match")
	}
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
