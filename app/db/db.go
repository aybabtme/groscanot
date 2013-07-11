package db

import (
	"github.com/aybabtme/dskvs"
	"github.com/robfig/revel"
	"log"
	"time"
)

var conf revel.MergedConfig
var Db *dskvs.Store

var KeySep = dskvs.CollKeySep

func init() {
	revel.OnAppStart(openDskvs)
}

func openDskvs() {
	dbPath, ok := revel.Config.String("dskvs.path")
	log.Printf("Init start")
	if !ok {
		panic("dskvs has no path to load!!!")
	}
	log.Printf("Loading dskvs on path '%s'", dbPath)
	t0 := time.Now()
	s, err := dskvs.Open(dbPath)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded in %s", time.Since(t0))
	Db = s
	tryRead()
}

func tryRead() {
	_, ok, err := Db.Get("topics/CEG")
	if err != nil {
		log.Printf("Got an error trying a read, %v", err)
		panic(err)
	}
	if !ok {
		log.Printf("Key was not found")
		return
	}
}
