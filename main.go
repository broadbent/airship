package main

import (
	"flag"
	"log"
	"time"

	"github.com/broadbent/airship/auctioneer"
	"gopkg.in/mgo.v2"
)

var mongoURL = "localhost:27017"
var interval, _ = time.ParseDuration("10m")

func main() {
	flag.Parse()

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	go auctioneer.Serve(session)
	auctioneer.Ticker(interval, session)

}
