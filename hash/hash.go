package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"text/template"
)

const (
	HashMethodMD5    = "md5"
	HashMethodSHA1   = "sha1"
	HashMethodSHA256 = "sha256"
	HashMethodSHA512 = "sha512"
)

var AllHashMethods = []string{
	HashMethodMD5,
	HashMethodSHA1,
	HashMethodSHA256,
	HashMethodSHA512,
}

func NewHash(method string, mac ...string) (hash.Hash, error) {
	var h func() hash.Hash
	switch method {
	case HashMethodMD5:
		h = md5.New
	case HashMethodSHA1:
		h = sha1.New
	case HashMethodSHA256:
		h = sha256.New
	case HashMethodSHA512:
		h = sha512.New
	default:
		return nil, fmt.Errorf("not support hash method: %s ", method)
	}
	if len(mac) == 1 {
		key := mac[0]
		if len(key) > 0 {
			return hmac.New(h, []byte(key)), nil
		}
	}
	return h(), nil
}

type Hasher struct {
	h     hash.Hash
	files []string

	ioOut io.Writer
	ioErr io.Writer
	tpl   *template.Template
}

func RunHasher(option *hashOption) error {
	hr := &Hasher{}
	var err error

	hr.h, err = NewHash(option.Method, option.HMac)
	if err != nil {
		return err
	}

	hr.ioErr = os.Stderr
	if option.OutputFile == "" {
		hr.ioOut = os.Stdout
	} else {
		os.Remove(option.OutputFile)
		f, err := os.OpenFile(option.OutputFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		hr.ioOut = f
		defer f.Close()
	}

	hr.tpl, err = template.New("").Parse(option.OutputFormat)
	if err != nil {
		return err
	}
	if option.EntireDir {
		hr.files, err = GetFiles(option.Path, option.Excludes)
		if err != nil {
			return err
		}
	} else {
		hr.files = []string{option.Path}
	}

	return hr.Run(option)
}

func (h *Hasher) Run(option *hashOption) error {
	if len(h.files) == 0 {
		return fmt.Errorf("no found any files")
	}

	buf := new(bytes.Buffer)
	for _, v := range h.files {
		buf.Reset()
		f, err := os.Open(v)
		if err != nil {
			h.LogError("open %s failed, %s", v, err.Error())
			continue
		}

		h.h.Reset()
		size, err := io.Copy(h.h, f)
		f.Close()
		if err != nil {
			h.LogError("do hash %s failed, %s", v, err.Error())
			continue
		}
		hbs := h.h.Sum(nil)
		hi := NewHashInfo(v, size, option.HumanReadable, hbs, option.ToUpper)
		err = h.tpl.Execute(buf, hi)
		if err != nil {
			h.LogError("template error: %s ", err.Error())
			continue
		}
		h.WriteLine(buf.Bytes())
	}
	return nil
}

func (h *Hasher) Hash(r io.Reader) ([]byte, error) {
	h.h.Reset()
	_, err := io.Copy(h.h, r)
	if err != nil {
		return nil, err
	}

	return h.h.Sum(nil), nil
}

func (h *Hasher) LogError(format string, args ...interface{}) {
	fmt.Fprintf(h.ioErr, format+"\n", args...)
}

func (h *Hasher) WriteLine(data []byte) error {
	_, err := h.ioOut.Write(data)
	h.ioOut.Write([]byte("\n"))
	return err
}
