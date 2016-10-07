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

	"github.com/broadbent/airship/config"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
)

var configuration *config.Configuration

var collectionNames = map[string]string{
	"user":    "user",
	"auction": "auction",
	"bid":     "bid",
}

func Ticker(session *mgo.Session, reset chan bool, newConfiguration *config.Configuration) {
	configuration = newConfiguration
	interval, _ := time.ParseDuration(configuration.Interval)
	ticker := time.NewTicker(interval)
	workers := make(chan bool, 1)
	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

	createAuction(session) //call for initial auction

	for {
		select {
		case <-ticker.C:
			go runWorker(workers, session)
			log.Println("Auction staging process started.")
		case <-workers:
			log.Println("Auction staging process complete.")
		case <-reset:
			ticker.Stop()
			ticker = time.NewTicker(interval)
			dropDatabase(session)
			createAuction(session)
			log.Println("Auction staging process reset.")
		case <-death:
			return
		}

	}
}

func writeResult(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if e := json.NewEncoder(w).Encode(obj); e != nil {
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

func reset(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	a.reset <- true
	return 200, nil
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
