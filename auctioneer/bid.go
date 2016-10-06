package auctioneer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

var bidIncrement = 1

type Bid struct {
	ID          string    `json:"id"`
	AuctionID   string    `json:"auction_id"`
	ItemID      string    `json:"item_id"`
	UserID      string    `json:"user_id"`
	Amount      int       `json:"amount"`
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

func findWinner(bids []Bid, winningBid Bid) Bid {
	for _, bid := range bids {
		fmt.Println(bid.Amount)
		var winningAmount = winningBid.Amount + bidIncrement
		fmt.Println(winningAmount)
		if bid.Amount > (winningBid.Amount + bidIncrement) {
			winningBid = bid
			fmt.Println("We have ourselves a new winner!")
		}
	}
	return winningBid
}

func findItemIndex(auction *Auction, bidID string) int {
	for index, item := range auction.Items {
		if item.ID == bidID {
			return index
		}
	}
	return -1
}

func placeBid(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var bid Bid

	bid.ID = xid.New().String()
	bid.TimeCreated = time.Now()

	body := readRequestBody(r)

	err := json.Unmarshal(body, &bid)
	if err != nil {
		panic(err)
	}

	col := a.session.DB(databaseName).C(collectionNames["auction"])

	query := bson.M{"id": bid.AuctionID, "live": true}

	fmt.Println(bid.AuctionID)

	auction := Auction{}
	err = col.Find(query).One(&auction)
	if err != nil {
		panic(err)
	}

	itemIndex := findItemIndex(&auction, bid.ItemID)
	bids := auction.Items[itemIndex].Bids
	bids = append(bids, bid)

	winner := findWinner(bids, auction.Items[itemIndex].Winning)
	fmt.Println(winner)

	if winner == bid {
		auction.Items[itemIndex].Winning = winner
		auction.Items[itemIndex].Bids = bids //should we anonmysie ids?
		query = bson.M{"id": auction.ID}
		err = col.Update(query, auction) //Won't be atomic - need to seperate out items, nodes, bids
		if err != nil {
			panic(err)
		}
	}

	writeResult(w, winner)

	return 200, nil
}
