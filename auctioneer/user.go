package auctioneer

import (
	"encoding/json"
	"log"

	"github.com/rs/xid"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type User struct {
	Name    string `json:"name,omitempty"`
	ID      string `json:"id,omitempty"`
	Balance int    `json:"balance,omitempty"`
}

func ensureUserIndex(s *mgo.Session) {
	c := s.DB(databaseName).C(collectionNames["user"])

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

func unmarshalUser(body []byte, user User) User {
	if e := json.Unmarshal(body, &user); e != nil {
		log.Panic(e)
	}
	return user
}

func addUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User
	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	col := a.session.DB(databaseName).C(collectionNames["user"])

	user.ID = xid.New().String()

	if e := col.Insert(user); e != nil {
		log.Panic(e)
	}

	writeResult(w, user)
	return 200, nil
}

func removeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(databaseName).C(collectionNames["user"])

	user := bson.M{"id": c.URLParams["user_id"]}

	if e := col.Remove(user); e != nil {
		log.Panic(e)
	}

	return 200, nil
}

func changeBalance(a *appContext, r *http.Request, deduct bool) {
	var user User
	var update bson.M

	col := a.session.DB(databaseName).C(collectionNames["user"])

	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	if deduct {
		update = bson.M{"$inc": bson.M{"balance": -user.Balance}}
	} else {
		update = bson.M{"$inc": bson.M{"balance": user.Balance}}
	}

	query := bson.M{"id": user.ID}

	if e := col.Update(query, update); e != nil {
		log.Panic(e)
	}

}

func addBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	changeBalance(a, r, false)

	return 200, nil
}

func deductBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	changeBalance(a, r, true)

	return 200, nil
}

func describeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(databaseName).C(collectionNames["user"])

	user := bson.M{"id": c.URLParams["user_id"]}
	result := User{}

	if e := col.Find(user).One(&result); e != nil {
		log.Panic(e)
	}

	writeResult(w, result)

	return 200, nil
}

func authenticateUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 501, nil
}
