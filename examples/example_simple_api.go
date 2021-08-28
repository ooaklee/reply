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
	Status       *reply.TransferObjectStatus `json:"status,omitempty"`
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

func (t *fooReplyTransferObject) SetStatus(transferObjectStatus *reply.TransferObjectStatus) {
	t.Bar.Status = transferObjectStatus
}

////////////////////

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var baseManifest []reply.ErrorManifest = []reply.ErrorManifest{
	{"example-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
}

var replier *reply.Replier = reply.NewReplier(baseManifest)

var replierWithCustomTransitionObj *reply.Replier = reply.NewReplier(baseManifest, reply.WithTransferObject(&fooReplyTransferObject{}))

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

	// mock passing additional headers in request
	mockAdditionalHeaders := map[string]string{
		"correlation-id": "d7c09ac2-fa46-4ece-bcde-1d7ad81d2230",
	}

	replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer:  w,
		Error:   unregisterdErr,
		Headers: mockAdditionalHeaders,
	})
}

func simpleTokensAPIHandler(w http.ResponseWriter, r *http.Request) {

	mockedAccessToken := "05e42c11-8bdd-423d-a2c1-c3c5c6604a30"
	mockedRefreshToken := "0e95c426-d373-41a5-bfe1-08db322527bd"

	replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer:       w,
		AccessToken:  mockedAccessToken,
		RefreshToken: mockedRefreshToken,
	})
}

func simpleAPIDefaultResponseHandler(w http.ResponseWriter, r *http.Request) {

	// Do something that only needs an empty response body, and 200 status code
	replier.NewHTTPResponse(&reply.NewResponseRequest{
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

func handleRequest() {
	var port string = ":8081"

	http.HandleFunc("/users", simpleUsersAPIHandler)
	http.HandleFunc("/users/3", simpleUsersAPINotFoundHandler)
	http.HandleFunc("/users/4", simpleUsersAPINoManifestEntryHandler)
	http.HandleFunc("/tokens/refresh", simpleTokensAPIHandler)
	http.HandleFunc("/defaults/1", simpleAPIDefaultResponseHandler)
	http.HandleFunc("/custom/users/3", simpleUsersAPINotFoundCustomReplierHandler)

	log.Printf("Serving simple API on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	handleRequest()
}
