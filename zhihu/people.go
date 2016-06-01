package main

import (
// "fmt"
// "io/ioutil"
// "os"

// "github.com/PuerkitoBio/goquery"
// "gopkg.in/mgo.v2"
)

type People struct {
	Name      string `bson:"name"`
	Location  string `bson:"location,omitempty"`
	WorkInfo  string `bson:"work_info"`
	Education string `bson:"edu"`
}

func inset() {
	// db := *mgo.Database
	// _ = db
}
