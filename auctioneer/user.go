package auctioneer

import (
	// "encoding/json"
	"fmt"

	"github.com/zenazn/goji/web"
	"net/http"
	// "gopkg.in/mgo.v2"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func addUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func removeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Println(string(c.URLParams["userId"]))
	return 200, nil
}

func addBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func deductBalance(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func describeUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func authenticateUser(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
