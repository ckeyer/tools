/*
 × 系统调用
**/

package main

import (
	"container/list"
	"github.com/howeyc/fsnotify"
	logpkg "log"
	"os"
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
	blk.ReadLogFile()
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
				if ev.IsModify() {
					log.Println("event:", ev)
					blk.ReadLogFile()
					log.Println("Add Over")
				}
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
