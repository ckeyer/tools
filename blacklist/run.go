/*
 × 系统调用
**/

package main

import (
	"container/list"
	"github.com/howeyc/fsnotify"
	"log"
)

const (
	watchFile     = "testfile.log"
	blacklistFile = "blacklist.txt"
	hostsDenyFile = "hosts.deny"
)

var (
	app *App
)

type App struct {
	A list.List
}

func init() {
	A = list.New()
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

	err = watcher.Watch(watchFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	watcher.Close()
}
