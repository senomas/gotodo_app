package service

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"time"
)

type Migration struct {
	Timestamp time.Time
	Filename  string
	Hash      string
	Result    string
	ID        int64
	Success   bool
}

func FileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filename, err)
	}
	defer f.Close()
	hasher := sha512.New()
	p := make([]byte, 1024)
	for {
		n, err := f.Read(p)
		hasher.Write(p[:n])
		if err == io.EOF {
			return fmt.Sprintf("%x", hasher.Sum(nil)), nil
		} else if err != nil {
			return "", err
		}
	}
}
