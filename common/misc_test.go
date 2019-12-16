package common

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"
	"reflect"
	"testing"
)

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
