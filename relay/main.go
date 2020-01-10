package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/Samyoul/storj-file-sender/common"
)

type stream struct {
	checksum []byte
	filename []byte
	sendConn chan net.Conn
	wg       sync.WaitGroup
}

func (s *stream) Close() {
	close(s.sendConn)
}

type streamMap map[string]*stream

func main() {
	// get init argument
	args := os.Args
	err := validateArgs(args)
	if err != nil {
		log.Fatalf("error - validating arguments : %s", err)
	}

	// open TCP server
	l, err := net.Listen("tcp", args[1])
	if err != nil {
		log.Fatalf("error - starting tcp server : %s", err)
	}
	defer l.Close()

	// Create streams map
	sm := streamMap{}

	// wait for connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return // return don't exit because you don't want to kill your whole server over a single fail connection
		}

		// serve connections with successful connection
		go handle(&sm, conn)
	}
}

func handle(sm *streamMap, conn net.Conn) {
	defer conn.Close()

	// Get the connection header
	// Added header so that I connection parameters can be exchanged between client and server
	hdr, err := common.GetRequestHeader(conn)
	if err != nil {
		log.Printf("error - getting request header - %s", err)
		return
	}

	// determine the type of connection coming in.
	// switch logic to give handler for each the send and receive requests
	switch string(hdr["Type"]) {
	case common.HeaderSend:
		err = send(sm, conn, hdr)
		if err != nil {
			log.Printf("error - processing send request - %s", err)
			return
		}
		break
	case common.HeaderReceive:
		err = receive(sm, conn, hdr)
		if err != nil {
			log.Printf("error - processing receive request - %s", err)
			return
		}
		break
	default:
		log.Println(hdr)
	}
}

func validateArgs(args []string) error {
	if len(args) != 2 {
		return errors.New(
			"invalid number of arguments.\n" +
			"expected : ./relay :<port>\n" +
			"example  : ./relay :9021")
	}

	return nil
}

func send(sm *streamMap, conn net.Conn, hdr common.Header) error {
	s := &stream{}
	(*sm)[string(hdr["Code"])] = s

	s.filename = hdr["Filename"]
	s.checksum = hdr["Checksum"]
	s.sendConn = make(chan net.Conn)
	s.wg = sync.WaitGroup{}
	defer s.Close()

	err := conn.(*net.TCPConn).SetReadBuffer(common.BufferLimit)
	if err != nil {
		return err
	}

	s.wg.Add(1)
	s.sendConn <- conn
	s.wg.Wait()

	/*
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		return err
	}
	*/

	delete(*sm, string(hdr["Code"]))

	return nil
}

func receive(sm *streamMap, conn net.Conn, hdr common.Header) error {
	// Check the stream exists in the stream map
	s, ok := (*sm)[string(hdr["Code"])]
	if !ok {
		return errors.New(fmt.Sprintf("unrecognised secret code '%s'", hdr["Code"]))
	}

	rh := common.MakeResponseHeaderReceive(s.filename, s.checksum)
	_, err := conn.Write(rh)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, <-s.sendConn)
	if err != nil {
		return err
	}

	s.wg.Done()

	return nil
}

// Most of this function is lifted straight from io.Copy
// however copying directly from one connection to another requires additional checks
// as a read on an open connection will never throw an error and therefore loop infinitely until the conn is closed.
// This has the affect of copying conns via the io.copy just hanging as neither conns will be closed
// So I've implemented a deliminator to check if the body transmission is complete
/*
func connCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	size := 32 * 1024
	buf := make([]byte, size)
	done := false

	for {
		nr, er := src.Read(buf)

		// Are the last two bytes the predefined data terminator
		if bytes.Compare(buf[nr-2:nr], common.Terminator) == 0 {
			nr = nr - 2
			done = true
		}
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])

			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if done {
			break
		}
	}
	return written, err
}
*/
