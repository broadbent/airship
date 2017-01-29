package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/broadbent/airship/auctioneer"
	"github.com/rs/xid"
)

var auctioneerRoot = "http://127.0.0.1:8080"
var userIDRoot = "service_"
var imageModifer = "_video"
var userID string
var phases = 3

// var bidIncrement = 10
// var startBid = 10
// var startPause = "10s"

var stagePause = "5s"
var phasePause = "7s"
var bidPause = "0.5s"
var minPause = 5
var maxPause = 15

var bidders = 1

var locations = map[string]int{
	"datacenter": 22,
	"residence":  6,
	"exchange":   6,
}

func main() {
	var userIDSuffix = flag.String("user", "a", "user ID suffice (a, b, c, etc.)")
	flag.Parse()
	userID = userIDRoot + *userIDSuffix

	startBidder()
}

func startBidder() bool {
	log.Printf("Bidder %v started.\n", userID)
	randomiseLocationQuotas()
	log.Printf("Quota is a follows: %v", locations)
	auctions := fetchAuctions(true)
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
		locations[location] = randomNumber(0, locations[location])
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
	auctions := fetchAuctions(true)
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
	auctions := fetchAuctions(false)
	for {
		if len(auctions) > 0 { //TODO: Only provisions the first auction
			items := winningAuctions(&auctions)
			provision(items)
			log.Printf("Ending provisioning phase.\n")
			return
		}
		time.Sleep(10 * time.Second)
		auctions = fetchAuctions(false)
	}
}

func provision(items []auctioneer.Item) {
	for _, item := range items {
		provisionItem(item)
	}
}

func provisionItem(item auctioneer.Item) {
	var provision auctioneer.Provision

	provision.Nodes = []string{item.ParentNode.ID}
	provision.ImageName = "lyndon160/" + userID + imageModifer
	provision.Memory = item.Memory
	provision.Hours = 1
	provision.PortBindings = make(map[string]int)
	provision.PortBindings["internal"] = 80

	makePost(provision, "/provision")
}

func winningAuctions(auctions *[]auctioneer.Auction) []auctioneer.Item {
	var items []auctioneer.Item

	for _, auction := range *auctions { //TODO: check if auction is no longer live
		for _, item := range auction.Items {
			if item.Leading.UserID == userID {
				items = append(items, item)
			}
		}
	}

	return items
}

func fetchAuctions(live bool) []auctioneer.Auction {
	var auctions []auctioneer.Auction

	var path string
	if live {
		path = auctioneerRoot + "/auction/live"
	} else {
		path = auctioneerRoot + "/auction/expired"
	}

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
		randomisedItems := randomiseTargetItems(auction.Items)
		for _, item := range randomisedItems {
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

func randomiseTargetItems(items []auctioneer.Item) []auctioneer.Item {
	dest := make([]auctioneer.Item, len(items))
	perm := rand.Perm(len(items))
	for i, v := range perm {
		dest[v] = items[i]
	}
	return dest
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
		bid.Valuation = bid.Valuation + randomNumber(10, 20)
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

	path := "/auction/bid"

	makePost(bid, path)
}

func makePost(obj interface{}, path string) {
	post, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(post)

	path = auctioneerRoot + path

	resp, err := http.Post(path, "application/json; charset=UTF-8", reader)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
