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
	ID       string `json:"id"`
	Memory   int    `json:"memory"`
	Location string `json:"location"`
	Bids     []Bid
	Winning  Bid
}

type Auction struct {
	ID          string    `json:"id"`
	Stage       int       `json:"stage"`
	TimeCreated time.Time `json:"time_created"`
	Live        bool      `json:"arch"`
	Items       []Item
	Nodes       []Node
}

func ensureAuctionIndex(s *mgo.Session) {
	c := s.DB(databaseName).C(collectionNames["auction"])

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
	col := a.session.DB(databaseName).C(collectionNames["auction"])

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
		auction.Items = append(auction.Items, createItems(&node)...)
	}
}

func createItems(node *Node) []Item {
	var items []Item

	slices := node.AvailableMemory / memorySplit
	for i := 0; i < slices; i++ {
		items = append(items, createItem(node.Location, memorySplit))
	}

	return items
}

func createItem(location string, memory int) Item {
	var item Item

	item.ID = xid.New().String()
	item.Location = location
	item.Memory = memory

	return item
}

func createAuction(s *mgo.Session) {
	var auction Auction
	var nodes []Node

	auction.Stage = 1
	auction.ID = xid.New().String()
	auction.TimeCreated = time.Now()
	auction.Live = true

	resp, err := http.Get(provisionerPath + "/nodes")
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

	col := s.DB(databaseName).C(collectionNames["auction"])

	if e := col.Insert(auction); e != nil {
		log.Panic(e)
	}

}

func updateAuction(s *mgo.Session, query bson.M, update bson.M) {
	col := s.DB(databaseName).C(collectionNames["auction"])

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
	query := bson.M{"stage": finalStage}
	update := bson.M{"$set": bson.M{"live": false}}

	updateAuction(s, query, update)
}
