package internal

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// Md5 returns the Md5 of a file passed as absolute path.
func Md5(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
