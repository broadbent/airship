package auctioneer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	ID              string `json:"id"`
	TotalMemory     int    `json:"total_memory"`
	AvailableMemory int    `json:"available_memory"`
	ReservedMemory  int    `json:"reserved_memory"`
	Location        string `json:"location"`
	Arch            string `json:"arch"`
}

type Item struct {
	ID              string `json:"id"`
	Memory          int    `json:"memory"`
	ParentNode      Node   `json:"parent_node"`
	ParentAuctionID string `json:"parent_auction_id"`
	Bids            []Bid  `json:"bids"`
	Leading         Bid    `json:"leading_bid"`
	Price           int    `json:"price"`
}

type Auction struct {
	ID    string    `json:"id"`
	Stage int       `json:"stage"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Live  bool      `json:"live"`
	Items []Item    `json:"items"`
	Nodes []Node    `json:"nodes"`
}

func ensureAuctionIndex(s *mgo.Session) {
	c := s.DB(configuration.DatabaseName).C(collectionNames["auction"])

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if e := c.EnsureIndex(index); e != nil {
		log.Panic(e)
	}
}

func listWonAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}

func listLostAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}

func listAuction(a *appContext, query bson.M) []Auction {
	col := a.session.DB(configuration.DatabaseName).C(collectionNames["auction"])

	var auctions []Auction
	if e := col.Find(query).All(&auctions); e != nil {
		log.Panic(e)
	}

	return auctions
}

func listLiveAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	query := bson.M{"live": true}

	writeResult(w, listAuction(a, query))

	return 200, nil
}

func listExpiredAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	query := bson.M{"live": false}

	writeResult(w, listAuction(a, query))

	return 200, nil
}

func describeAuction(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	query := bson.M{"id": c.URLParams["auction_id"]}

	writeResult(w, listAuction(a, query))

	return 200, nil
}

func sliceNodes(auction *Auction) {
	for _, node := range auction.Nodes {
		auction.Items = append(auction.Items, createItems(node, auction.ID)...)
	}
}

func createItems(node Node, auctionID string) []Item {
	var items []Item

	slices := node.AvailableMemory / configuration.MemorySplit
	for i := 0; i < slices; i++ {
		items = append(items, createItem(node, auctionID, configuration.MemorySplit))
	}

	return items
}

func createStartingBid(auctionID string, itemID string) Bid {
	bid := createBid()

	bid.AuctionID = auctionID
	bid.ItemID = itemID
	bid.Valuation = configuration.StartingValuation
	//TODO: Set bid.UserID to an admin user

	return bid
}

func createItem(node Node, auctionID string, memory int) Item {
	var item Item

	item.ID = xid.New().String()
	item.ParentNode = node
	item.ParentAuctionID = auctionID
	item.Memory = memory
	item.Leading = createStartingBid(item.ID, auctionID)

	return item
}

func calcStopTime(start time.Time) time.Time {
	interval, _ := time.ParseDuration(configuration.Interval)
	stages := configuration.FinalStage - 1
	duration := interval * time.Duration(stages)
	stop := start.Add(duration)
	return stop
}

func createAuction(s *mgo.Session) {
	var auction Auction
	var nodes []Node

	auction.Stage = 1
	auction.ID = xid.New().String()
	auction.Start = time.Now()
	auction.End = calcStopTime(auction.Start)
	auction.Live = true

	resp, err := http.Get(configuration.ProvisionerPath + "/nodes")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Panic(err)
	}

	auction.Nodes = nodes

	sliceNodes(&auction)

	col := s.DB(configuration.DatabaseName).C(collectionNames["auction"])

	if e := col.Insert(auction); e != nil {
		log.Panic(e)
	}

}

func updateAuction(s *mgo.Session, query bson.M, update bson.M) {
	col := s.DB(configuration.DatabaseName).C(collectionNames["auction"])

	_, err := col.UpdateAll(query, update)
	if err != nil {
		log.Panic(err)
	}
}

func transitionAuctionStage(s *mgo.Session) {
	query := bson.M{}
	update := bson.M{"$inc": bson.M{"stage": 1}}

	updateAuction(s, query, update)
}

func expireAuctions(s *mgo.Session) {
	query := bson.M{"stage": configuration.FinalStage}
	update := bson.M{"$set": bson.M{"live": false}}

	updateAuction(s, query, update)
}
