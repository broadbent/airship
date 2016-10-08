package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/broadbent/airship/auctioneer"
)

var auctioneerRoot = "http://localhost:8080"

func main() {
	auctions := discoverAuctions()
	items := determineTargetItems(auctions, []string{"datacenter"})
	log.Println(len(items))
	bids := generateBids(items)
	bid(bids)
	provision()
}

func discoverAuctions() []auctioneer.Auction {
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

	log.Println(string(body))

	return auctions

}

func determineTargetItems(auctions []auctioneer.Auction, locationCriterium []string) []auctioneer.Item {
	var items []auctioneer.Item

	for _, auction := range auctions {
		for _, item := range auction.Items {
			for _, locationCriteria := range locationCriterium {
				if item.ParentNode.Location == locationCriteria {
					items = append(items, item)
				}
			}
		}
	}

	return items
}

func generateBids(items []auctioneer.Item) []auctioneer.Bid {
	var bids []auctioneer.Bid

	for _, item := range items {
		var bid auctioneer.Bid

		bid.AuctionID = item.ParentAuctionID
		bid.ItemID = item.ID
		bid.UserID = "test"
		bid.Valuation = 10

		bids = append(bids, bid)
	}

	return bids
}

func bid(bids []auctioneer.Bid) {
	for _, bid := range bids {
		executeBid(bid)
	}
}

func executeBid(bid auctioneer.Bid) {
	log.Print(bid)
	post, err := json.Marshal(bid)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(post)

	path := auctioneerRoot + "/auction/bid"
	log.Println(path)
	resp, err := http.Post(path, "application/json; charset=UTF-8", reader)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(body))
}
func provision() {

}
