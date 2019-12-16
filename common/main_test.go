package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMakeRequestHeaderSend(t *testing.T) {
	filename := "rad-music.mp4"
	code := "abounding-metallic-gold-cockatiel-218"
	checksum := []byte{
		0xea, 0x7b, 0xe7, 0x3b, 0x1f, 0x65, 0xc2, 0x5a,
		0x7e, 0x45, 0x51, 0x6b, 0xec, 0xc9, 0xf3, 0x3a,
		0x75, 0x6a, 0x3c, 0xc6, 0x27, 0xee, 0x36, 0x6b,
		0x07, 0x69, 0x70, 0xcd, 0xc4, 0xe2, 0xef, 0x44,
	}

	hdr := MakeRequestHeaderSend(filename, code, checksum)
	spew.Dump(hdr)

	hr, err := GetRequestHeader(bytes.NewReader(hdr))
	spew.Dump(hr, err)
}

func Test_HashFile(t *testing.T) {
	cs := []struct {
		Filename  string
		Checksum  string
		ExpectHsr hash.Hash
		ExpectErr error
		FailMsg   string
	}{
		{"../.gitignore", "ea7be73b1f65c25a7e45516becc9f33a756a3cc627ee366b076970cdc4e2ef44", sha256.New(), nil, "Existing file - Expects checksum of : '' and no error. Got '%s'"},
		{"../.gitignored", "", nil, new(os.PathError), "Existing file - an error. Got '%s'"},
	}

	for _, c := range cs {
		hsr, err := HashFile(c.Filename)
		if reflect.TypeOf(err) != reflect.TypeOf(c.ExpectErr) {
			t.Errorf(c.FailMsg, err)
			continue
		}

		if reflect.TypeOf(hsr) != reflect.TypeOf(c.ExpectHsr) {
			t.Errorf(c.FailMsg, err)
			continue
		}

		if hsr == nil {
			continue
		}

		checksum := hex.EncodeToString(hsr.Sum(nil))
		if checksum != c.Checksum {
			t.Errorf(c.FailMsg, checksum)
		}
	}

}
