package auctioneer

import (
	"log"
	"net/http"

	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
)

type contextHandler func(*appContext, web.C, http.ResponseWriter, *http.Request) (int, error)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc contextHandler
}

type Routes []Route

type appContext struct {
	session *mgo.Session
	reset   chan bool
}

type appHandler struct {
	*appContext
	H contextHandler
}

func (ah appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	status, err := ah.H(ah.appContext, c, w, r)

	if err != nil {
		log.Println("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func Serve(s *mgo.Session, reset chan bool) {

	dropDatabase(s)
	ensureUserIndex(s)
	ensureAuctionIndex(s)

	context := &appContext{session: s, reset: reset}

	r := web.New()

	for _, route := range routes {
		switch route.Method {
		case "Get", "get", "GET":
			r.Get(route.Pattern, appHandler{context, route.HandlerFunc})
		case "Post", "post", "POST":
			r.Post(route.Pattern, appHandler{context, route.HandlerFunc})
		}
	}

	graceful.ListenAndServe(":8080", r)
}

var routes = Routes{
	Route{
		"AuthenticateUser",
		"GET",
		"/user/authenticate",
		authenticateUser,
	},
	Route{
		"AddUser",
		"POST",
		"/user/add",
		addUser,
	},
	Route{
		"RemoveUser",
		"GET",
		"/user/remove/:user_id",
		removeUser,
	},
	Route{
		"AddBalance",
		"POST",
		"/user/balance/add",
		addBalance,
	},
	Route{
		"DeductBalance",
		"POST",
		"/user/balance/deduct",
		deductBalance,
	},
	Route{
		"DescribeUser",
		"GET",
		"/user/:user_id",
		describeUser,
	},
	Route{
		"ListBids",
		"GET",
		"/user/bid/all/:user_id",
		listBids,
	},
	Route{
		"ListAcceptedBids",
		"GET",
		"/user/bid/accepted/:user_id",
		listAcceptedBids,
	},
	Route{
		"ListRejectedBids",
		"GET",
		"/user/bid/rejected/:user_id",
		listRejectedBids,
	},
	Route{
		"ListWonAuctions",
		"GET",
		"/user/auction/won/:user_id",
		listWonAuctions,
	},
	Route{
		"ListLostAuctions",
		"GET",
		"/user/auction/lost/:user_id",
		listLostAuctions,
	},
	Route{
		"ListLiveAuctions",
		"GET",
		"/auction/live",
		listLiveAuctions,
	},
	Route{
		"ListExpiredAuctions",
		"GET",
		"/auction/expired",
		listExpiredAuctions,
	},
	Route{
		"DescribeAuction",
		"GET",
		"/auction/:auction_id",
		describeAuction,
	},
	Route{
		"Index",
		"POST",
		"/auction/bid",
		placeBid,
	},
	Route{
		"Provision",
		"POST",
		"/provision",
		provision,
	},
	Route{
		"Reset",
		"GET",
		"/debug/reset",
		reset,
	},
}
