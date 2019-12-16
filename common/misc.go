package common

import (
	"hash"
	"os"
	"io"
	"crypto/sha256"
)

const (
	Kb          = 1024
	Mb          = Kb * Kb
	BufferLimit = 32 * Kb
)

var (
	Terminator = []byte{0xBC, 0x00}
)

func HashFile(fn string) (hash.Hash, error) {
	h := sha256.New()

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h, nil
}
