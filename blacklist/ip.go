package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type IPAddr struct {
	Id uint32
	IP [4]byte
}

func (this *IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", this.IP[0], this.IP[1], this.IP[2], this.IP[3])
}

func NewIPAddr(s string) *IPAddr {
	if len(s) < 7 {
		return nil
	}
	var p [4]byte
	var id uint32 = 0
	// s = strings.Trim(s, "\r")
	bs := strings.Split(s, ".")
	if len(bs) != 4 {
		panic(errors.New("IP Length Error :#" + s + "#"))
	}
	for i, v := range bs {
		n, e := strconv.Atoi(v)
		if e != nil {
			panic(e)
		}
		id += uint32(n) << uint32((3-i)*8)
		p[i] = byte(n)
	}
	return &IPAddr{
		Id: id,
		IP: p,
	}
}

func (this *IPAddr) Equal(b *IPAddr) bool {
	if this.Id == b.Id {
		return true
	}
	return false
}
