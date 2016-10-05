package auctioneer

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var database = "airship"
var collection = "user"

type User struct {
	Name    string `json:"name,omitempty"`
	UserId  string `json:"userid,omitempty"`
	Balance int    `json:"balance,omitempty"`
}

func ensureUserIndex(s *mgo.Session) {
	c := s.DB("airship").C("user")

	index := mgo.Index{
		Key:        []string{"userid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func readRequestBody(r *http.Request) (body []byte) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	return body
}

func unmarshalUser(body []byte, user User) User {
	err := json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}
	return user
}

func writeResult(w http.ResponseWriter, user User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		panic(err)
	}

}

func addUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User
	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	col := a.session.DB(database).C(collection)

	user.UserId = xid.New().String()

	err := col.Insert(user)
	if err != nil {
		panic(err)
	}

	writeResult(w, user)
	return 200, nil
}

func removeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(database).C(collection)

	user := bson.M{"userid": c.URLParams["userId"]}

	err := col.Remove(user)
	if err != nil {
		panic(err)
	}

	return 200, nil
}

func changeBalance(col *mgo.Collection, user User, balance int) {
	query := bson.M{"userid": user.UserId}
	change := bson.M{"$inc": bson.M{"balance": balance}}

	err := col.Update(query, change)
	if err != nil {
		panic(err)
	}

}

func addBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User

	col := a.session.DB(database).C(collection)

	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	changeBalance(col, user, user.Balance)

	return 200, nil
}

func deductBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User

	col := a.session.DB(database).C(collection)

	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	changeBalance(col, user, -user.Balance)

	return 200, nil
}

func describeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(database).C(collection)

	user := bson.M{"userid": c.URLParams["userId"]}

	result := User{}
	err := col.Find(user).One(&result)
	if err != nil {
		panic(err)
	}

	writeResult(w, result)

	return 200, nil
}

func authenticateUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}
