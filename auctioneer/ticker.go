package auctioneer

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zenazn/goji/web"
)

func Ticker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	workers := make(chan bool, 1)
	death := make(chan os.Signal, 1)
	signal.Notify(death, os.Interrupt, os.Kill)

	for {
		select {
		case <-ticker.C:
			log.Println("Scheduled task is triggered.")
			go runWorker(workers)
		case <-workers:
			log.Println("Scheduled task is completed.")
		case <-death:
			return
		}
	}
}

func runWorker(workers chan bool) bool {
	log.Println("working")
	workers <- true
	return true
}

func provision(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
