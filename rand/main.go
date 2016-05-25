package main

import (
	"flag"
	"fmt"
	"github.com/ckeyer/commons/lib"
	"os"
	"strconv"
	"strings"
)

var (
	rtype string

	start = flag.Int("start", 0, "set start of rand number")
	end   = flag.Int("end", 100, "set end of rand number")

	randMap = map[string]randHandle{
		"s": func(l int) string {
			fmt.Println("...")
			return lib.RandomString(l)
		},
		"n": func(l int) string {
			ns := make([]string, 0, l)
			for i := 0; i < l; i++ {
				ns = append(ns, fmt.Sprint(lib.RandomInt(*start, *end)))
			}
			return strings.Join(ns, ", ")
		},
	}
)

type randHandle func(l int) string

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("error args length...")
		os.Exit(1)
	}
	if l, err := strconv.Atoi(args[1]); err == nil && l > 0 {
		if rf, ok := randMap[args[0]]; ok {
			fmt.Println("out: ", rf(l))
		} else {
			fmt.Println("error, not found ", args[0])
		}
	} else {
		fmt.Println(args[1], "not a number")
	}
}

func main() {
	fmt.Println()
}
