package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var (
	showAll  = flag.Bool("a", false, "show all files")
	forHuman = flag.Bool("h", false, "show for human")
	showLong = flag.Bool("l", false, "show long info")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println(34116906298 / 1024 / 1024 / 8)
	return
	args := flag.Args()
	for _, v := range args {
		fmt.Printf("%s: \n", v)
		s, err := GetSize(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
		return
		f, err := os.Open(v)
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		inf, err := f.Stat()
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		if inf.IsDir() {
			fmt.Println("is dir")
		}
		fmt.Println(inf.ModTime())
		fmt.Println(inf.Mode())
		fmt.Println(inf.Size())
		if stat, ok := inf.Sys().(*syscall.Stat_t); ok {
			fmt.Printf("%+v\n", stat)
		}
	}
}

func GetSize(fpath string) (int64, error) {
	// fmt.Println("fpath:", fpath)
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		fmt.Println("error:", err.Error())
		return -1, err
	}
	inf, err := f.Stat()
	if err != nil {
		fmt.Println("error:", err.Error())
		return -1, err
	}
	if inf.IsDir() {
		return getDirSize(fpath)
	}
	return inf.Size(), nil
}

func getDirSize(dir string) (int64, error) {
	var count int64
	err := filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || err != nil {
			fmt.Println("errorrrr", err)
			os.Exit(0)
			return err
		}
		dpath, _ := filepath.Abs(dir)
		if dpath == fpath {
			return nil
		}

		size, err := GetSize(fpath)
		if err != nil {
			return nil
		}
		count += size
		return nil
	})
	// fmt.Println("count: ", count)
	return count, err
}
