package auctioneer

import (
	"fmt"
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
	db *mgo.Session
}

type appHandler struct {
	*appContext
	H contextHandler
}

func (ah appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	status, err := ah.H(ah.appContext, c, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
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

func Serve(s *mgo.Session) {

	context := &appContext{db: s}

	r := web.New()

	for _, route := range routes {
		switch route.Method {
		case "Get", "get", "GET":
			r.Get(route.Pattern, appHandler{context, route.HandlerFunc})
			fmt.Println("GET")
		case "Post", "post", "POST":
			r.Post(route.Pattern, appHandler{context, route.HandlerFunc})
			fmt.Println("POST")
		}
		fmt.Println(route.Pattern)
	}

	graceful.ListenAndServe(":8080", r)

	// 	router.
	// 		Methods(route.Method).
	// 		Path(route.Pattern).
	// 		Name(route.Name).
	// 		Handler(handler)
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
		"/user/remove/:userId",
		removeUser,
	},
	Route{
		"AddBalance",
		"POST",
		"/user/balance/add/{userId}",
		addBalance,
	},
	Route{
		"DeductBalance",
		"POST",
		"/user/balance/deduct/{userId}",
		deductBalance,
	},
	Route{
		"DescribeUser",
		"GET",
		"/user/{userId}",
		describeUser,
	},
	Route{
		"ListBids",
		"GET",
		"/user/bid/all/{userId}",
		listBids,
	},
	Route{
		"ListAcceptedBids",
		"GET",
		"/user/bid/accepted/{userId}",
		listAcceptedBids,
	},
	Route{
		"ListRejectedBids",
		"GET",
		"/user/bid/rejected/{userId}",
		listRejectedBids,
	},
	Route{
		"ListWonAuctions",
		"GET",
		"/user/auction/won/{userId}",
		listWonAuctions,
	},
	Route{
		"ListLostAuctions",
		"GET",
		"/user/auction/lost/{userId}",
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
		"/auction/describe/{auctionId}",
		describeAuction,
	},
	Route{
		"Index",
		"POST",
		"/auction/bid/{auctionId}",
		placeBid,
	},
	Route{
		"Provision",
		"POST",
		"/auction/provision/{auctiontId}",
		provision,
	},
}
