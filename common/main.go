package common

import (
	"errors"
	"fmt"
	"io"
	"bytes"
)

const (
	Kb          = 1024
	Mb          = Kb * Kb
	BufferLimit = 32 * Kb

	HeaderSend    = "s"
	HeaderReceive = "r"
)

var (
	Terminator = []byte{0xBC, 0x00}
)

type Header map[string][]byte

func MakeRequestHeaderSend(file string, code string, checksum []byte) []byte {
	hdr := *new([]byte)

	hdr = append(hdr, HeaderSend...)

	secret := make([]byte, 64)
	copy(secret[:], code)
	hdr = append(hdr, secret...)

	cs := make([]byte, 32)
	copy(cs[:], checksum)
	hdr = append(hdr, cs...)

	hdr = append(hdr, file...)
	hdr = append(hdr, Terminator...)

	return hdr
}

func MakeRequestHeaderReceive(code string) []byte {
	hdr := *new([]byte)

	hdr = append(hdr, HeaderReceive...)

	secret := make([]byte, 64)
	copy(secret[:], code)
	hdr = append(hdr, secret...)

	return hdr
}

func GetRequestHeader(conn io.Reader) (Header, error) {
	hdr := Header{}

	// Get request type
	t := make([]byte, 1)
	if _, err := conn.Read(t); err != nil {
		return hdr, err
	}
	hdr["Type"] = t

	// Get request code
	c := make([]byte, 64)
	if _, err := conn.Read(c); err != nil {
		return hdr, err
	}
	hdr["Code"] = c

	switch string(t) {
	case HeaderSend:
		return getRequestHeaderSend(hdr, conn)
	case HeaderReceive:
		return getRequestHeaderReceive(hdr, conn)
	default:
		return hdr, errors.New(fmt.Sprintf("request type not recognised '%s'", t))
	}
}

func getRequestHeaderSend(hdr Header, conn io.Reader) (Header, error) {
	var err error

	// Get request checksum
	cs := make([]byte, 32)
	if _, err = conn.Read(cs); err != nil {
		return hdr, err
	}
	hdr["Checksum"] = cs

	// Get request filename
	if hdr["Filename"], err = getFileName(conn); err != nil {
		return hdr, err
	}

	return hdr, nil
}

func getRequestHeaderReceive(hdr Header, conn io.Reader) (Header, error) {
	return hdr, nil
}

func getFileName(conn io.Reader) ([]byte, error) {
	var out []byte

	buff := make([]byte, 1)
	for {
		if _, err := conn.Read(buff); err != nil {
			return nil, err
		}
		out = append(out, buff...)

		if bytes.Compare(Terminator[:1], buff) == 0 {

			if _, err := conn.Read(buff); err != nil {
				return nil, err
			}
			out = append(out, buff...)

			if bytes.Compare(Terminator[1:], buff) == 0 {
				break
			}
		}
	}

	return out[:len(out)-2], nil
}

func MakeResponseHeaderReceive(file []byte, checksum []byte) []byte {
	hdr := *new([]byte)

	cs := make([]byte, 32)
	copy(cs[:], checksum)
	hdr = append(hdr, cs...)

	hdr = append(hdr, file...)
	hdr = append(hdr, Terminator...)

	return hdr
}

func GetResponseHeader(conn io.Reader) (Header, error) {
	hdr := Header{}
	var err error

	// Get request checksum
	cs := make([]byte, 32)
	if _, err = conn.Read(cs); err != nil {
		return hdr, err
	}
	hdr["Checksum"] = cs

	// Get request filename
	if hdr["Filename"], err = getFileName(conn); err != nil {
		return hdr, err
	}

	return hdr, nil
}
