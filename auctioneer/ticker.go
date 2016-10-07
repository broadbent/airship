package auctioneer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gopkg.in/mgo.v2"
)

var bidIncrement = 1
var memorySplit = 256
var finalStage = 2
var provisionerPath = "http://148.88.226.119:60000"

var databaseName = "airship"
var collectionNames = map[string]string{
	"user":    "user",
	"auction": "auction",
	"bid":     "bid",
}

func Ticker(interval time.Duration, session *mgo.Session) {
	ticker := time.NewTicker(interval)
	workers := make(chan bool, 1)
	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

	createAuction(session) //call for initial auction

	for {
		select {
		case <-ticker.C:
			log.Println("Auction staging process started.")
			go runWorker(workers, session)
		case <-workers:
			log.Println("Auction staging process complete.")
		case <-death:
			return
		}
	}
}

func writeResult(w http.ResponseWriter, user interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if e := json.NewEncoder(w).Encode(user); e != nil {
		log.Panic(e)
	}
}

func readRequestBody(r *http.Request) (body []byte) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	if e := r.Body.Close(); e != nil {
		log.Panic(e)
	}

	return body
}

func dropDatabase(s *mgo.Session) {
	if e := s.DB("airship").DropDatabase(); e != nil {
		log.Panic(e)
	}
}

func runWorker(workers chan bool, s *mgo.Session) bool {
	transitionAuctionStage(s)
	expireAuctions(s)
	createAuction(s)
	workers <- true
	return true
}
