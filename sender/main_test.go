package main

import (
	"testing"
	"errors"
	"reflect"
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
