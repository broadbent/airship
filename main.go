package main

import (
	"flag"
	"time"

	"github.com/broadbent/airship/auctioneer"
	"gopkg.in/mgo.v2"
)

var mongoURL = "localhost:27017"
var interval, _ = time.ParseDuration("60s")

func main() {
	flag.Parse()

	session, _ := mgo.Dial(mongoURL)

	go auctioneer.Serve(session)
	auctioneer.Ticker(interval)

}
