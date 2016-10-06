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

var databaseName = "airship"
var collectionNames = map[string]string{
	"user":    "user",
	"auction": "auction",
	"bid":     "bid",
}
var provisionerPath = "http://148.88.226.119:60000"

func Ticker(interval time.Duration, session *mgo.Session) {
	ticker := time.NewTicker(interval)
	workers := make(chan bool, 1)
	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

	createAuction(session) //call for initial auction

	for {
		select {
		case <-ticker.C:
			log.Println("Scheduled task is triggered.")
			go runWorker(workers, session)
		case <-workers:
			log.Println("Scheduled task is completed.")
		case <-death:
			return
		}
	}
}

func writeResult(w http.ResponseWriter, user interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		panic(err)
	}
}

func readRequestBody(r *http.Request) (body []byte) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	return body
}

func dropDatabase(s *mgo.Session) {
	err := s.DB("airship").DropDatabase()
	if err != nil {
		panic(err)
	}
}

func runWorker(workers chan bool, s *mgo.Session) bool {
	transitionAuctionStage(s)
	expireAuctions(s)
	createAuction(s)
	workers <- true
	return true
}
