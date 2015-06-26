package main

import (
	// "errors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type IBlockList interface {
	InitOldList()
	Add(ip string)
	Operate()
}

type IPAddrs []*IPAddr

type BlackList struct {
	OldIps IPAddrs
	NewIps IPAddrs
}

func (this IPAddrs) Len() int {
	return len(this)
}
func (this IPAddrs) Less(i, j int) bool {
	if this[i].Id < this[j].Id {
		return true
	}
	return false
}
func (this IPAddrs) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *BlackList) InitOldList() {
	fout, err := os.OpenFile(blacklistFile, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	fd, err := ioutil.ReadAll(fout)

	items := strings.Split(string(fd), "\r\n")
	for _, v := range items {
		blk.AddOld(v)
	}
}
func (this *BlackList) Add(ip string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in blacklist.go Add ", err)
		}
	}()
	item := NewIPAddr(ip)
	if item == nil {
		// err := errors.New("Parse IP Error " + ip)
		return
	}
	this.addNewItem(item)
}
func (this *BlackList) AddOld(ip string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in blacklist.go AddOld ", err)
		}
	}()
	item := NewIPAddr(ip)
	if item == nil {
		return
	}
	for _, v := range this.OldIps {
		if v.Equal(item) {
			return
		}
	}
	this.addToOld(item)
}
func (this *BlackList) addNewItem(ip *IPAddr) {
	for _, v := range this.OldIps {
		if v.Equal(ip) {
			return
		}
	}
	for _, v := range this.NewIps {
		if v.Equal(ip) {
			return
		}
	}
	this.addToNew(ip)
}
func (this *BlackList) addToNew(ip *IPAddr) {
	this.NewIps = append(this.NewIps, ip)
	sort.Sort(blk.NewIps)
}
func (this *BlackList) addToOld(ip *IPAddr) {
	this.OldIps = append(this.OldIps, ip)
	sort.Sort(blk.OldIps)
}

func (this *BlackList) WriteTxt() {
	fout, err := os.OpenFile(blacklistFile, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	for _, v := range this.NewIps {
		fout.WriteString(v.String() + "\r\n")
	}
}

func (this *BlackList) WriteDeny() {
	fout, err := os.OpenFile(hostsDenyFile, os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	for _, v := range this.NewIps {
		fout.WriteString("ALL: " + v.String() + "\r\n")
	}
	this.NewIps = nil
}
func (this *BlockList) Operate() {
	this.WriteDeny()
	this.WriteTxt()
}
