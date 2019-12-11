package codegen

import (
	"testing"
	"github.com/davecgh/go-spew/spew"
)

func Test_Make(t *testing.T) {

	for i:=0; i<100; i++{
		spew.Dump(Make(int64(i)))
	}
}