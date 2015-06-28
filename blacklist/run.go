/*
 × 系统调用
**/

package main

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	logFile       = "testfile.log"
	blacklistFile = "blacklist.txt"
	hostsDenyFile = "hosts.deny"
	indexString   = "Failed password for root from"
)

var (
	blk BlackList
)

func init() {
	blk.InitOldList()
	sl := ReadLogFile()
	items := strings.Split(sl, "\n")
	for _, v := range items {
		if strings.Contains(v, indexString) {
			handleLog(v)
		}
	}
	blk.WriteDeny()
	blk.WriteTxt()
	for _, v := range blk.TmpIps {
		log.Println(v)
	}
}
func handleLog(s string) {
	for i, v := range strings.Fields(s) {
		if i == 10 {
			blk.Add(v)
		}
	}
}
func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)

	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(logFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	watcher.Close()
}

func ReadLogFile() string {
	fi, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
