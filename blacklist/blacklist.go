package main

import (
	"errors"
	// "log"
	"sort"
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

func (this *BlackList) Add(ip string) {
	item := NewIPAddr(ip)
	if item == nil {
		err := errors.New("Parse IP Error " + ip)
		panic(err)
	}
	this.addItem(item)
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
