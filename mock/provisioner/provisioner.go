package main

import (
    "fmt"
    "net/http"
    "log"
)

// PathPrinter prints the URL path called by airship
func PathPrinter(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Println("call", r.URL.Path)
}

// main creates a server to listen for calls from airship.
func main() {
    http.HandleFunc("/", PathPrinter)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal(err)
    }
}