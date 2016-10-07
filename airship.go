package main

import (
	"flag"
	"log"

	"github.com/broadbent/airship/auctioneer"
	"github.com/broadbent/airship/config"
	"gopkg.in/mgo.v2"
)

func main() {
	var path = flag.String("path", "config.json", "path for configuration file")
	flag.Parse()

	configuration := config.Read(*path)

	session := createMongoDBSession(configuration.MongoURI)
	defer session.Close()

	reset := make(chan bool, 1)

	go auctioneer.Serve(session, reset)
	auctioneer.Ticker(session, reset, &configuration)

}

func createMongoDBSession(mongoURI string) *mgo.Session {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.Panic(err)
	}

	return session
}
