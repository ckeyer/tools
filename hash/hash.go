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
	"runtime"
	"sync"
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
	id            int
	humanReadable bool
	toUpper       bool

	h     hash.Hash
	ioErr io.Writer
	tpl   *template.Template

	exitCh chan int
}

func RunHasher(option *hashOption) error {
	routineNum := runtime.NumCPU() * 2
	filesBuf := make(chan string, routineNum*2)
	writeBuf := make(chan string, routineNum)
	closeAll := make(chan int)
	var err error
	files := []string{}
	if option.EntireDir {
		files, err = GetFiles(option.Path, option.Excludes)
		if err != nil {
			return err
		}
	} else {
		files = []string{option.Path}
	}

	var ioOut io.Writer
	if option.OutputFile == "" {
		ioOut = os.Stdout
	} else {
		os.Remove(option.OutputFile)
		f, err := os.OpenFile(option.OutputFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		ioOut = f
		defer f.Close()
	}

	var wg sync.WaitGroup
	go func() {
		wg.Add(len(files))
		for _, file := range files {
			select {
			case filesBuf <- file:
			case <-closeAll:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case data := <-writeBuf:
				wg.Done()
				ioOut.Write([]byte(data))
			case <-closeAll:
				return
			}
		}
	}()

	for i := 0; i < routineNum; i++ {
		hr := &Hasher{
			id:            i + 1,
			humanReadable: option.HumanReadable,
			toUpper:       option.ToUpper,
			ioErr:         os.Stderr,
			exitCh:        closeAll,
		}
		hr.h, err = NewHash(option.Method, option.HMac)
		if err != nil {
			return err
		}
		hr.tpl, err = template.New("").Parse(option.OutputFormat)
		if err != nil {
			return err
		}
		go hr.Run(filesBuf, writeBuf)
	}

	wg.Wait()
	close(closeAll)
	return nil
}

func (h *Hasher) Run(fileCh chan string, out chan string) {
	buf := new(bytes.Buffer)
	for {
		select {
		case name := <-fileCh:
			// time.Sleep(time.Second)
			buf.Reset()
			f, err := os.Open(name)
			if err != nil {
				h.LogError("open %s failed, %s", name, err.Error())
				continue
			}
			h.h.Reset()
			size, err := io.Copy(h.h, f)
			f.Close()
			if err != nil {
				h.LogError("do hash %s failed, %s", name, err.Error())
				continue
			}
			hbs := h.h.Sum(nil)
			hi := NewHashInfo(name, size, h.humanReadable, hbs, h.toUpper)
			err = h.tpl.Execute(buf, hi)
			if err != nil {
				h.LogError("template error: %s ", err.Error())
				continue
			}
			out <- buf.String() + "\n"
		case <-h.exitCh:
			return
		}
	}
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
