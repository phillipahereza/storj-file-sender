package common

import (
	"errors"
	"fmt"
	"io"
)

const (
	Kb          = 1024
	Mb          = Kb * Kb
	BufferLimit = 32 * Kb

	HeaderSend    = "s"
	HeaderReceive = "r"
)

type Header map[string][]byte

func MakeRequestHeaderSend(code string, checksum []byte) []byte {
	hdr := *new([]byte)

	hdr = append(hdr, HeaderSend...)

	secret := make([]byte, 64)
	copy(secret[:], code)
	hdr = append(hdr, secret...)

	cs := make([]byte, 32)
	copy(cs[:], checksum)
	hdr = append(hdr, cs...)

	return hdr
}

func GetRequestHeader(conn io.Reader) (Header, error) {
	hdr := Header{}

	t := make([]byte, 1)
	if _, err := conn.Read(t); err != nil {
		return hdr, err
	}
	hdr["Type"] = t

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

	c := make([]byte, 64)
	if _, err := conn.Read(c); err != nil {
		return hdr, err
	}
	hdr["Code"] = c

	cs := make([]byte, 32)
	if _, err := conn.Read(cs); err != nil {
		return hdr, err
	}
	hdr["Checksum"] = cs

	return hdr, nil
}

func getRequestHeaderReceive(hdr Header, conn io.Reader) (Header, error) {
	// TODO

	return hdr, nil
}
