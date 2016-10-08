package auctioneer

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

type Bid struct {
	ID          string    `json:"id"`
	AuctionID   string    `json:"auction_id"`
	ItemID      string    `json:"item_id"`
	UserID      string    `json:"user_id"`
	Valuation   int       `json:"valuation"`
	TimeCreated time.Time `json:"time_created"`
}

func listAcceptedBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}

func listRejectedBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}

func listBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}

func isLeader(bid Bid, price int) bool { //needs new logic
	if bid.Valuation > (price + configuration.BidIncrement) {
		return true
		log.Println("We have ourselves a new winner!")
	}
	return false

}

func findItemIndex(auction *Auction, bidID string) int {
	for index, item := range auction.Items {
		if item.ID == bidID {
			return index
		}
	}
	return -1
}

func createBid() Bid {
	var bid Bid

	bid.ID = xid.New().String()
	bid.TimeCreated = time.Now()

	return bid
}

func placeBid(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	bid := createBid()

	body := readRequestBody(r)

	err := json.Unmarshal(body, &bid)
	if err != nil {
		panic(err)
	}

	col := a.session.DB(configuration.DatabaseName).C(collectionNames["auction"])

	query := bson.M{"id": bid.AuctionID, "live": true}
	auction := Auction{}

	err = col.Find(query).One(&auction)
	if err != nil {
		panic(err)
	}

	itemIndex := findItemIndex(&auction, bid.ItemID)
	item := &auction.Items[itemIndex]

	bids := append(item.Bids, bid)

	leading := isLeader(bid, item.Leading.Valuation)

	change := false

	if leading {
		item.Price = item.Leading.Valuation
		item.Leading = bid
		change = true
	} else if bid.Valuation > item.Price {
		item.Price = bid.Valuation
		change = true
	}

	if change {
		item.Bids = bids //should we anonmysie ids?
		query = bson.M{"id": auction.ID}
		err = col.Update(query, auction) //Won't be atomic - need to seperate out items, nodes, bids
		if err != nil {
			panic(err)
		}
	}

	writeResult(w, item)

	return 200, nil
}
