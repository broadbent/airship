package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/broadbent/airship/auctioneer"
	"github.com/rs/xid"
)

var auctioneerRoot = "http://localhost:8080"
var userIDRoot = "bidder_"
var userID string
var phases = 3

// var bidIncrement = 10
// var startBid = 10
// var startPause = "10s"
var stagePause = "5s"
var phasePause = "10s"
var bidPause = "2s"
var minPause = 10
var maxPause = 20

var bidders = 1

var locations = map[string]int{
	"datacenter": 4,
	"residence":  4,
}

func main() {
	startBidder(randomNumber(1, 1000))
}

func startBidder(bidderNumber int) bool {
	userID = userIDRoot + strconv.Itoa(bidderNumber)
	log.Printf("Bidder %v started.\n", bidderNumber)
	randomiseLocationQuotas()
	log.Println(locations)
	auctions := fetchAuctions()
	items := determineTargetItems(auctions, locations)
	bids := generateBids(items)
	artificialSleep("", true)
	biddingStage(bids)
	artificialSleep(stagePause, false)
	provisioningStage()

	return true
}

func randomiseLocationQuotas() {
	for location, _ := range locations {
		locations[location] = randomNumber(0, 5)
	}
}

func artificialSleep(duration string, random bool) {
	var sleep time.Duration

	if random {
		sleep = randomDuration(minPause, maxPause)
	} else {
		sleep, _ = time.ParseDuration(duration)
	}

	time.Sleep(sleep)

}

func randomDuration(min, max int) time.Duration {
	duration := time.Duration(randomNumber(min, max)) * time.Second
	return duration
}

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min

}

func biddingStage(bids map[string]*auctioneer.Bid) {
	for i := 0; i < phases; i++ {
		biddingPhase(bids, i)
	}
}

func biddingPhase(bids map[string]*auctioneer.Bid, phaseNumber int) {
	log.Printf("Starting bidding phase %v.\n", phaseNumber)
	executeBids(bids)
	printBids(bids)
	incrementBids(bids)
	printLeading()
	log.Printf("Ending bidding phase %v.\n", phaseNumber)
	artificialSleep(phasePause, false)
}

func printBids(bids map[string]*auctioneer.Bid) {
	log.Println("Bid summary is as follows:")
	for _, bid := range bids {
		log.Printf("Bid for %v, valued at %v, was submitted.", bid.UserTag, bid.Valuation)
	}
}

func printLeading() { //could also check against user_tag?
	log.Println("Leading the following auctions:")
	auctions := fetchAuctions()
	for _, auction := range auctions {
		for _, item := range auction.Items {
			if item.Leading.UserID == userID {
				log.Printf("Leading item %v.", item.ID)
			}
		}
	}
}

func provisioningStage() {
	log.Printf("Starting provisioning phase.\n")
	provision()
	// check if finished
	//call provision in won items
	log.Printf("Ending provisioning phase.\n")
}

func fetchAuctions() []auctioneer.Auction {
	var auctions []auctioneer.Auction

	path := auctioneerRoot + "/auction/live"

	resp, err := http.Get(path) //should the end point not be '/provision_docker_containers'?
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &auctions)

	return auctions

}

func determineTargetItems(auctions []auctioneer.Auction, locationCriterium map[string]int) []auctioneer.Item {
	var items []auctioneer.Item

	for _, auction := range auctions {
		for _, item := range auction.Items {
			for location, quota := range locationCriterium {
				if item.ParentNode.Location == location {
					if quota > 0 {
						items = append(items, item)
						locationCriterium[location] = quota - 1
					}
				}
			}
		}
	}

	return items
}

func generateBids(items []auctioneer.Item) map[string]*auctioneer.Bid {
	bids := make(map[string]*auctioneer.Bid)

	for _, item := range items {
		var bid auctioneer.Bid

		bid.UserTag = xid.New().String()
		bid.AuctionID = item.ParentAuctionID
		bid.ItemID = item.ID
		bid.UserID = userID
		bid.Valuation = randomNumber(10, 20)

		bids[bid.UserTag] = &bid
	}

	return bids
}

func incrementBids(bids map[string]*auctioneer.Bid) {
	for _, bid := range bids {
		bid.Valuation += randomNumber(1, 5)
	}
}

func executeBids(bids map[string]*auctioneer.Bid) {
	for _, bid := range bids {
		artificialSleep(bidPause, false)
		executeBid(bid)
	}
}

func executeBid(bid *auctioneer.Bid) {
	log.Printf("Placing bid now: %v.\n", bid.UserTag)

	post, err := json.Marshal(bid)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(post)

	path := auctioneerRoot + "/auction/bid"

	resp, err := http.Post(path, "application/json; charset=UTF-8", reader)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// var item auctioneer.Item

	// json.Unmarshal(body, &item)

}
func provision() {

}
