package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/broadbent/airship/db"
)

var dbFile = flag.String("db", "/tmp/airship.db", "Path to the BoltDB file")
var buckets = []string{"api", "transaction", "resource"}

func main() {
	flag.Parse()

	var db db.DB
	if err := db.Open(*dbFile, 0600, buckets); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	timestamp := db.Write("log", "", "/test/test", true)

}
