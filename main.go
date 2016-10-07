package main

import (
	"flag"
	"log"
	"time"

	"github.com/broadbent/airship/auctioneer"
	"gopkg.in/mgo.v2"
)

var mongoURL = "localhost:27017"
var interval, _ = time.ParseDuration("5s")

func main() {
	flag.Parse()

	reset := make(chan bool, 1)
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	go auctioneer.Serve(session, reset)
	auctioneer.Ticker(interval, session, reset)

}
