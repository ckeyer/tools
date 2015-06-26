/*
 × 系统调用
**/

package main

import (
	"container/list"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	logpkg "log"
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
	app *App
	log *logpkg.Logger
)

type App struct {
	A list.List
}

func init() {
	log = logpkg.New(os.Stdout, "# BLK #: ", logpkg.Lshortfile)

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
	for _, v := range blk.NewIps {
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
