package auctioneer

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func listWonAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func listLostAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func listLiveAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func listExpiredAuctions(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func describeAuction(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
