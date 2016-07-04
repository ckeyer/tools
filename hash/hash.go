package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ckeyer/commons/fileutils"
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

type Hasher struct {
	h     hash.Hash
	files []string

	infos map[string]*HashInfo
	// formater *Formater
	stdOut io.Writer
	stdErr io.Writer
	tpl    *template.Template
}

func NewHasher(option *hashOption) (*Hasher, error) {
	hr := &Hasher{}
	hr.infos = make(map[string]*HashInfo)
	var err error

	hr.h, err = NewHash(option.Method, option.HMac)
	if err != nil {
		return nil, err
	}

	hr.stdErr = os.Stderr
	if option.OutputFile == "" {
		hr.stdOut = os.Stdout
	} else {
		f, err := os.OpenFile(option.OutputFile, os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return nil, err
		}
		hr.stdOut = f
	}

	tpl, err := template.New("").Parse(option.OutputFormat)
	if err != nil {
		return nil, err
	}
	hr.tpl = tpl

	return hr, nil
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
		return nil, fmt.Errorf("not support hash method: %s", method)
	}
	if len(mac) == 1 {
		key := mac[0]
		if len(key) > 0 {
			return hmac.New(h, []byte(key)), nil
		}
	}
	return h(), nil
}

func GetFiles(path string, excludes []string) ([]string, error) {
	includes := []string{}

	err := filepath.Walk(path, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || err != nil {
			return err
		}

		matched, err := fileutils.Matches(strings.TrimPrefix(fpath, path), excludes)
		if err != nil || matched {
			return err
		}

		if f.Mode().IsRegular() {
			fmt.Println("include: ", strings.TrimPrefix(fpath, path))
			includes = append(includes, fpath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return includes, nil
}

func (h *Hasher) Hash(r io.Reader) ([]byte, error) {
	h.h.Reset()
	_, err := io.Copy(h.h, r)
	if err != nil {
		return nil, err
	}

	return h.h.Sum(nil), nil
}
