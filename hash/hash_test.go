package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestHashFile(t *testing.T) {
	return
	filename := "tests/a.txt"
	// filename := "./hash"
	f, err := os.Open(filename)
	if err != nil {
		t.Errorf("open file %s failed, error: %s", filename, err.Error())
		t.FailNow()
	}
	defer f.Close()

	h, err := NewHash("md5")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
		return
	}

	io.Copy(h, f)
	hash := h.Sum(nil)
	t.Logf("file: %s, hash: %x ", filename, hash)

}

func TestGetFiles(t *testing.T) {
	fs, err := GetFiles("../", []string{".git"})

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
		return
	}

	buf := new(bytes.Buffer)
	for _, path := range fs {
		fi, _ := os.Stat(path)
		size := fi.Size()

		f, _ := os.Open(path)
		siz, _ := io.Copy(buf, f)
		f.Close()
		if size != siz {
			t.Errorf("%s not equae size: %d, siz: %d ", path, size, siz)
		}
	}
}
