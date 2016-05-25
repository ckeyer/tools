package main

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type IBlockList interface {
	InitOldList()
	Add(ip string) int
	Operate()
}

func InitBlockList(blk IBlockList) {
	blk.InitOldList()
}
func AddNew(blk IBlockList, ip string) {
	if blk.Add(ip) > 0 {
		blk.Operate()
	}
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

// 初始化：获取已有黑名单列表
func (this *BlackList) InitOldList() {
	fout, err := os.OpenFile(hostsDenyFile, os.O_RDONLY, 0666)
	if err == os.ErrNotExist {
		fout, err = os.OpenFile(hostsDenyFile, os.O_WRONLY|os.O_CREATE, 0666)
	}
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	fd, err := ioutil.ReadAll(fout)

	items := strings.Split(string(fd), "\n")
	for _, v := range items {
		if strings.HasPrefix(v, "#") {
			continue
		}
		if strings.TrimSpace(v) == "" {
			continue
		}

		kv := strings.SplitN(v, ":", 2)
		blk.AddOld(kv[1])
	}
}

// 读取日志文件
func (this *BlackList) ReadLogFile() {
	fi, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		fi.Close()
		if err := recover(); err != nil {
			log.Printf("RECOVER ReadLogFile: %v\n", err)
		}
	}()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	items := strings.Split(string(fd), "\n")
	for _, v := range items {
		if strings.Contains(v, indexString) {
			this.handleLog(v)
		}
	}
	this.WriteDeny()
	this.copyNew2Old()
	this.WriteTxt()
	this.clearNew()
}
func (this *BlackList) handleLog(s string) {
	for i, v := range strings.Fields(s) {
		if i == 10 {
			this.Add(v)
		}
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
	sort.Sort(this.NewIps)
}
func (this *BlackList) addToOld(ip *IPAddr) {
	this.OldIps = append(this.OldIps, ip)
	sort.Sort(this.OldIps)
}

func (this *BlackList) WriteTxt() {
	fout, err := os.OpenFile(blacklistFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	if this.NewIps != nil {
		for _, v := range this.OldIps {
			fout.WriteString("\n" + v.String())
			if err != nil {
				panic(err)
			}
		}
	}
}
func (this *BlackList) WriteDeny() {
	fout, err := os.OpenFile(hostsDenyFile, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	for _, v := range this.NewIps {
		_, err := fout.WriteString("\n" + "ALL: " + v.String())
		if err != nil {
			panic(err)
		}
	}
}
func (this *BlackList) Operate() {
	this.WriteDeny()
	this.WriteTxt()
}
func (this *BlackList) copyNew2Old() {
	for _, v := range this.NewIps {
		this.addToOld(v)
	}
}
func (this *BlackList) clearNew() {
	this.NewIps = nil
}
