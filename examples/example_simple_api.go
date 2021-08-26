package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/ooaklee/reply"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var baseManifest []reply.ErrorManifest = []reply.ErrorManifest{
	{"example-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
}

var replier *reply.Replier = reply.NewReplier(baseManifest)

func simpleUsersAPINotFoundHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errors.New("example-404-error")

	replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

func simpleUsersAPIHandler(w http.ResponseWriter, r *http.Request) {

	mockedQueriedUsers := []user{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Sam Smith"},
	}

	replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Data:   mockedQueriedUsers,
	})
}

func simpleUsersAPINoManifestEntryHandler(w http.ResponseWriter, r *http.Request) {

	// unregisterdErr an error that's  unregistered in manifest
	// should return 500
	unregisterdErr := errors.New("unexpected-error")

	replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  unregisterdErr,
	})
}

func handleRequest() {
	var port string = ":8081"

	http.HandleFunc("/users", simpleUsersAPIHandler)
	http.HandleFunc("/users/3", simpleUsersAPINotFoundHandler)
	http.HandleFunc("/users/4", simpleUsersAPINoManifestEntryHandler)

	log.Printf("Serving simple API on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	handleRequest()
}
