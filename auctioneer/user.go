package auctioneer

import (
	"encoding/json"

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

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func unmarshalUser(body []byte, user User) User {
	err := json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}
	return user
}

func addUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User
	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	col := a.session.DB(databaseName).C(collectionNames["user"])

	user.ID = xid.New().String()

	err := col.Insert(user)
	if err != nil {
		panic(err)
	}

	writeResult(w, user)
	return 200, nil
}

func removeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(databaseName).C(collectionNames["user"])

	user := bson.M{"id": c.URLParams["user_id"]}

	err := col.Remove(user)
	if err != nil {
		panic(err)
	}

	return 200, nil
}

func changeBalance(col *mgo.Collection, user User, balance int) {
	query := bson.M{"id": user.ID}
	change := bson.M{"$inc": bson.M{"balance": balance}}

	err := col.Update(query, change)
	if err != nil {
		panic(err)
	}

}

func addBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User

	col := a.session.DB(databaseName).C(collectionNames["user"])

	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	changeBalance(col, user, user.Balance)

	return 200, nil
}

func deductBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var user User

	col := a.session.DB(databaseName).C(collectionNames["user"])

	body := readRequestBody(r)
	user = unmarshalUser(body, user)

	changeBalance(col, user, -user.Balance)

	return 200, nil
}

func describeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	col := a.session.DB(databaseName).C(collectionNames["user"])

	user := bson.M{"id": c.URLParams["user_id"]}

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
