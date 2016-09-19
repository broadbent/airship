package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// main reads the list of actions and then enacts each in order.
func main() {
    actions := LoadActionFile("actions.txt")
    for _, action := range actions {
        a := strings.Split(string(action), ",")
        fmt.Println(a[0], a[1], Action(a[0], a[1]))
    }
}

// LoadActionFile loads action.txt and parses each line into an array.
func LoadActionFile(filename string) (actions []string) {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }
    return strings.Split(string(content), "\n")
}

// Action erforms a specific action, given the method name.
func Action(method string, params string) (outcome bool) {
    if method == "wait" {
        outcome = ActionWait(params)
    } else if method == "call" {
        outcome = ActionCall(params)
    }
    return outcome
}

// ActionWait realises the 'wait' method. Sleeps for a specific duration (in seconds).
func ActionWait(params string) bool {
    period, err := strconv.Atoi(params)
    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(time.Duration(period) * time.Second)
    return true
}

// ActionCall realises the 'call' method. Makes a HTTP GET request to the given URL.
// Returns the success based upon the HTTP status code returned from the request.
func ActionCall(params string) bool {
    code, _ := MakeHTTPCall(params)
    if (code >= 200) && (code < 300) {
        return true
    } 
    return false
}

// MakeHTTPCall makes the HTTP call to the given URL. It then reads the body, and returns
// this, along with the HTTP status code.
func MakeHTTPCall(url string) (int, string) {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    return resp.StatusCode, string(body)
}