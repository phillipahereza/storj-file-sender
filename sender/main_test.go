package main

import (
	"testing"
	"errors"
	"reflect"
	"github.com/davecgh/go-spew/spew"
	"os"
)

func Test_validateFlags(t *testing.T) {

	cs := []struct{
		HostFlag string
		FileFlag string
		ExpectErr error
		FailMsg string
	}{
		{"", "", errors.New(""), "Two blank inputs, expect an error"},
		{"localhost", "", errors.New(""), "Missing file, expect an error"},
		{"", "README.md", errors.New(""), "Missing host, expect an error"},
		{"localhost", "README.md", nil,"both inputs present, expect no error"},
	}

	for _, c := range cs {
		err := validateFlags(&c.HostFlag, &c.FileFlag)
		if reflect.TypeOf(err) != reflect.TypeOf(c.ExpectErr) {
			t.Error(c.FailMsg)
		}
	}
}

func Test_checksumFile(t *testing.T) {
	cs := []struct{
		Filename string
		Checksum string
		ExpectErr error
		FailMsg string
	}{
		{"../.gitignore", "ea7be73b1f65c25a7e45516becc9f33a756a3cc627ee366b076970cdc4e2ef44", nil, "Existing file - Expects checksum of : '' and no error. Got '%s'"},
		{"../.gitignored", "", new(os.PathError), "Existing file - an error. Got '%s'"},
	}

	for _, c := range cs {
		checksum, err := checksumFile(&c.Filename)
		if reflect.TypeOf(err) != reflect.TypeOf(c.ExpectErr) {

			spew.Dump(err, c.ExpectErr)
			t.Errorf(c.FailMsg, err)
		}

		if checksum != c.Checksum {
			t.Errorf(c.FailMsg, checksum)
		}
	}

}
