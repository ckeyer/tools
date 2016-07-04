package main

import (
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

	fi, _ := os.Stat(filename)
	var aaa io.Reader
	_ = aaa
	io.Copy(h, f)
	hash := h.Sum(nil)
	t.Logf("file: %s, hash: %x, size: %v", filename, hash, NewSize(fi.Size(), true))
	t.Error("...")
}

func TestGetFiles(t *testing.T) {
	fs, err := GetFiles("../", "../", []string{})

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
		return
	}

	t.Logf("%+v", fs)
	t.Error("...")
}
