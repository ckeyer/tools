package main

import (
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type IPAddrs []*IPAddr

type BlackList struct {
	OldIps IPAddrs
	NewIps IPAddrs
	TmpIps IPAddrs
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

	items := strings.Split(string(fd), "\n")
	for _, v := range items {
		blk.AddOld(v)
	}
}
func (this *BlackList) Add(ip string) {
	item := NewIPAddr(ip)
	if item == nil {
		err := errors.New("Parse IP Error " + ip)
		panic(err)
	}
	this.addItem(item)
}
func (this *BlackList) AddOld(ip string) {
	item := NewIPAddr(ip)
	if item == nil {
		err := errors.New("Parse IP Error " + ip)
		panic(err)
	}
	for _, v := range this.OldIps {
		if v.Equal(item) {
			return
		}
	}
	this.addToOld(item)
}
func (this *BlackList) addItem(ip *IPAddr) {
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
	for _, v := range this.TmpIps {
		if v.Equal(ip) {
			return
		}
	}
	this.addToTmp(ip)
}
func (this *BlackList) addToNew(ip *IPAddr) {
	this.NewIps = append(this.NewIps, ip)
	sort.Sort(blk.NewIps)
}
func (this *BlackList) addToOld(ip *IPAddr) {
	this.OldIps = append(this.OldIps, ip)
	sort.Sort(blk.OldIps)
}
func (this *BlackList) addToTmp(ip *IPAddr) {
	this.TmpIps = append(this.TmpIps, ip)
	sort.Sort(blk.TmpIps)
}

func (this *BlackList) WriteTxt() {
	fout, err := os.OpenFile(blacklistFile, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	for _, v := range this.TmpIps {
		fout.WriteString(v.String() + "\r\n")
	}
}

func (this *BlackList) WriteDeny() {
	fout, err := os.OpenFile(hostsDenyFile, os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	for _, v := range this.TmpIps {
		fout.WriteString("ALL: " + v.String() + "\r\n")
	}
}
