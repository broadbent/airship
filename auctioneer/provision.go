package auctioneer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Provision struct {
	Nodes     []string `json:"nodes,omitempty"`
	ImageName string   `json:"image_name"`
	Memory    int      `json:"ram"`
	Hours     int      `json:"hours"`
	UserID    string   `json:"user_id,omitempty"`
	AuctionID string   `json:"auction_id,omitempty"`
}

//[List nodes] [String image_name] [Int ram] [Int hours]

func checkValidity(col *mgo.Collection, auctionID string) {
	query := bson.M{"id": auctionID}

	auction := Auction{}
	err := col.Find(query).One(&auction)
	if err != nil {
		panic(err)
	}
}

func resolveNodes() []string {
	var nodes []string
	return nodes
}

func provision(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var provision Provision

	body := readRequestBody(r)

	err := json.Unmarshal(body, &provision)
	if err != nil {
		panic(err)
	}

	col := a.session.DB(databaseName).C(collectionNames["auction"])

	checkValidity(col, provision.AuctionID)

	provision.Nodes = resolveNodes()
	provision.UserID = ""
	provision.AuctionID = ""

	post, err := json.Marshal(provision)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(post)

	path := provisionerPath + "/provision_dockers"
	fmt.Println(path)
	resp, err := http.Post(path, "application/json; charset=UTF-8", reader) //should the end point not be '/provision_docker_containers'?
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	return 200, nil
}
