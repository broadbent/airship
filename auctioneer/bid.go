package auctioneer

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func listAcceptedBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func listRejectedBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func listBids(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func placeBid(a *appContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
