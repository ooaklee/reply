package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/ooaklee/reply"
)

/////////////////////////////////////////////////
/////// Custom Transition Object Example ////////
// This is an example of how you can create a
// custom response structure based on your
// requirements.

type fooReplyTransferObject struct {
	HTTPWriter http.ResponseWriter `json:"-"`
	Headers    map[string]string   `json:"-"`
	StatusCode int                 `json:"-"`
	Bar        barEmbeddedExample  `json:"bar,omitempty"`
}

type barEmbeddedExample struct {
	Errors       []reply.TransferObjectError `json:"errors,omitempty"`
	Meta         map[string]interface{}      `json:"meta,omitempty"`
	Data         interface{}                 `json:"data,omitempty"`
	AccessToken  string                      `json:"access_token,omitempty"`
	RefreshToken string                      `json:"refresh_token,omitempty"`
}

func (t *fooReplyTransferObject) SetHeaders(headers map[string]string) {
	t.Headers = headers
}

func (t *fooReplyTransferObject) SetStatusCode(code int) {
	t.StatusCode = code
}

func (t *fooReplyTransferObject) SetMeta(meta map[string]interface{}) {
	t.Bar.Meta = meta
}

func (t *fooReplyTransferObject) SetWriter(writer http.ResponseWriter) {
	t.HTTPWriter = writer
}

func (t *fooReplyTransferObject) SetAccessToken(token string) {
	t.Bar.AccessToken = token
}

func (t *fooReplyTransferObject) SetRefreshToken(token string) {
	t.Bar.RefreshToken = token
}

func (t *fooReplyTransferObject) GetWriter() http.ResponseWriter {
	return t.HTTPWriter
}

func (t *fooReplyTransferObject) GetStatusCode() int {
	return t.StatusCode
}

func (t *fooReplyTransferObject) SetData(data interface{}) {
	t.Bar.Data = data
}

func (t *fooReplyTransferObject) RefreshTransferObject() reply.TransferObject {
	return &fooReplyTransferObject{}
}

func (t *fooReplyTransferObject) SetErrors(transferObjectErrors []reply.TransferObjectError) {
	t.Bar.Errors = transferObjectErrors
}

////////////////////

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var baseManifest []reply.ErrorManifest = []reply.ErrorManifest{
	{"example-404-error": reply.ErrorManifestItem{Title: "resource not found", StatusCode: http.StatusNotFound}},
}

var replier *reply.Replier = reply.NewReplier(baseManifest)

var replierWithCustomTransitionObj *reply.Replier = reply.NewReplier(baseManifest, reply.WithTransferObject(&fooReplyTransferObject{}))

func simpleUsersAPINotFoundHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errors.New("example-404-error")

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

func simpleUsersAPIHandler(w http.ResponseWriter, r *http.Request) {

	mockedQueriedUsers := []user{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Sam Smith"},
	}

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Data:   mockedQueriedUsers,
	})
}

func simpleUsersAPINoManifestEntryHandler(w http.ResponseWriter, r *http.Request) {

	// unregisterdErr an error that's  unregistered in manifest
	// should return 500
	unregisterdErr := errors.New("unexpected-error")

	// mock passing additional headers in request
	mockAdditionalHeaders := map[string]string{
		"correlation-id": "d7c09ac2-fa46-4ece-bcde-1d7ad81d2230",
	}

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer:  w,
		Error:   unregisterdErr,
		Headers: mockAdditionalHeaders,
	})
}

func simpleTokensAPIHandler(w http.ResponseWriter, r *http.Request) {

	mockedAccessToken := "05e42c11-8bdd-423d-a2c1-c3c5c6604a30"
	mockedRefreshToken := "0e95c426-d373-41a5-bfe1-08db322527bd"

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer:       w,
		AccessToken:  mockedAccessToken,
		RefreshToken: mockedRefreshToken,
	})
}

func simpleAPIDefaultResponseHandler(w http.ResponseWriter, r *http.Request) {

	// Do something that only needs an empty response body, and 200 status code
	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
	})
}

func simpleUsersAPINotFoundCustomReplierHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errors.New("example-404-error")

	replierWithCustomTransitionObj.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

//////////////////////////////
//// Handlers Using Aides ////

func simpleUsersAPINotFoundUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errors.New("example-404-error")

	_ = replier.NewHTTPErrorResponse(w, serverErr)
}

func simpleUsersAPIUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	mockedQueriedUsers := []user{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Sam Smith"},
	}

	_ = replier.NewHTTPDataResponse(w, http.StatusCreated, mockedQueriedUsers)
}

func simpleUsersAPINoManifestEntryUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// unregisterdErr an error that's  unregistered in manifest
	// should return 500
	unregisterdErr := errors.New("unexpected-error")

	// mock passing additional headers in request
	mockAdditionalHeaders := map[string]string{
		"correlation-id": "d7c09ac2-fa46-4ece-bcde-1d7ad81d2230",
	}

	_ = replier.NewHTTPErrorResponse(w, unregisterdErr, reply.WithHeaders(mockAdditionalHeaders))
}

func simpleTokensAPIUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	mockedAccessToken := "05e42c11-8bdd-423d-a2c1-c3c5c6604a30"
	mockedRefreshToken := "0e95c426-d373-41a5-bfe1-08db322527bd"

	_ = replier.NewHTTPTokenResponse(w, http.StatusOK, mockedAccessToken, mockedRefreshToken)
}

func simpleAPIDefaultResponseUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// Do something that only needs an empty response body.
	// Note: 200 status code will be returned if status code passed is 0.
	// Otherwise passed code would be used
	_ = replier.NewHTTPBlankResponse(w, http.StatusOK)
}

func simpleUsersAPINotFoundCustomReplierUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errors.New("example-404-error")

	replierWithCustomTransitionObj.NewHTTPErrorResponse(w, serverErr)
}

/////////////////////////////

func handleRequest() {
	var port string = ":8081"

	http.HandleFunc("/users", simpleUsersAPIHandler)
	http.HandleFunc("/users/3", simpleUsersAPINotFoundHandler)
	http.HandleFunc("/users/4", simpleUsersAPINoManifestEntryHandler)
	http.HandleFunc("/tokens/refresh", simpleTokensAPIHandler)
	http.HandleFunc("/defaults/1", simpleAPIDefaultResponseHandler)
	http.HandleFunc("/custom/users/3", simpleUsersAPINotFoundCustomReplierHandler)

	http.HandleFunc("/aides/users", simpleUsersAPIUsingAideHandler)
	http.HandleFunc("/aides/users/3", simpleUsersAPINotFoundUsingAideHandler)
	http.HandleFunc("/aides/users/4", simpleUsersAPINoManifestEntryUsingAideHandler)
	http.HandleFunc("/aides/tokens/refresh", simpleTokensAPIUsingAideHandler)
	http.HandleFunc("/aides/defaults/1", simpleAPIDefaultResponseUsingAideHandler)
	http.HandleFunc("/aides/custom/users/3", simpleUsersAPINotFoundCustomReplierUsingAideHandler)

	log.Printf("Serving simple API on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	handleRequest()
}
