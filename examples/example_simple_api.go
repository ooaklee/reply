package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (t *fooReplyTransferObject) SetTokenOne(token string) {
	t.Bar.AccessToken = token
}

func (t *fooReplyTransferObject) SetTokenTwo(token string) {
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

/////////////////////////////////////////////////
//// Custom Transition Object Error Example /////
// This is an example of how you can create a
// custom response structure for errrors retutned
// in the response

type barError struct {

	// Title a short summary of the problem
	Title string `json:"title,omitempty"`

	// Message a description of the error
	Message string `json:"message,omitempty"`

	// About holds the link that gives further insight into the error
	About string `json:"about,omitempty"`

	// More randomd top level attribute to make error
	// difference
	More struct {
		// Status the HTTP status associated with error
		Status string `json:"status,omitempty"`

		// Code internal error code used to reference error
		Code string `json:"code,omitempty"`

		// Meta contains additional meta-information about the error
		Meta interface{} `json:"meta,omitempty"`
	} `json:"more,omitempty"`
}

// SetTitle adds title to error
func (b *barError) SetTitle(title string) {
	b.Title = title
}

// GetTitle returns error's title
func (b *barError) GetTitle() string {
	return b.Title
}

// SetDetail adds detail to error
func (b *barError) SetDetail(detail string) {
	b.Message = detail
}

// GetDetail return error's detail
func (b *barError) GetDetail() string {
	return b.Message
}

// SetAbout adds about to error
func (b *barError) SetAbout(about string) {
	b.About = about
}

// GetAbout return error's about
func (b *barError) GetAbout() string {
	return b.About
}

// SetStatusCode converts and add http status code to error
func (b *barError) SetStatusCode(status int) {
	b.More.Status = strconv.Itoa(status)
}

// GetStatusCode returns error's HTTP status code
func (b *barError) GetStatusCode() string {
	return b.More.Status
}

// SetCode adds internal code to error
func (b *barError) SetCode(code string) {
	b.More.Code = code
}

// GetCode returns error's internal code
func (b *barError) GetCode() string {
	return b.More.Code
}

// SetMeta adds meta property to error
func (b *barError) SetMeta(meta interface{}) {
	b.More.Meta = meta
}

// GetMeta returns error's meta property
func (b *barError) GetMeta() interface{} {
	return b.More.Meta
}

// RefreshTransferObject returns an empty instance of transfer object
// error
func (b *barError) RefreshTransferObject() reply.TransferObjectError {
	return &barError{}
}

////////////////////

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// handle errors
var (
	// errExample404 is an example 404 error
	errExample404 = errors.New("example-404-error")
	// errExampleDobValidation is an example dob validation error
	errExampleDobValidation = errors.New("example-dob-validation-error")
	// errExampleNameValidation is an example name validation error
	errExampleNameValidation = errors.New("example-name-validation-error")
	// errExampleMissing is an example missing error
	errExampleMissing = errors.New("example-missing-error")
)

// Example implementation of Error Manifest
var baseManifest []reply.ErrorManifest = []reply.ErrorManifest{
	{errExample404: reply.ErrorManifestItem{Title: "resource not found", StatusCode: http.StatusNotFound}},
	{errExampleNameValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"}},
	{errExampleDobValidation: reply.ErrorManifestItem{Title: "Validation Error", Detail: "Check your DoB, and try again.", Code: "100YT", StatusCode: http.StatusBadRequest}},
}

// Replier with default Transition Object & Transition Object Error
var replier *reply.Replier = reply.NewReplier(baseManifest)

// Replier with custom Transition Object & default Transition Object Error
var replierWithCustomTransitionObj *reply.Replier = reply.NewReplier(baseManifest, reply.WithTransferObject(&fooReplyTransferObject{}))

// Replier with standard Transition Object & custom Transition Object Error
var replierWithCustomTransitionObjError *reply.Replier = reply.NewReplier(baseManifest, reply.WithTransferObjectError(&barError{}))

// Replier with custom Transition Object & custom Transition Object Error
var replierWithCustomTransitionObjs *reply.Replier = reply.NewReplier(baseManifest, reply.WithTransferObjectError(&barError{}), reply.WithTransferObject(&fooReplyTransferObject{}))

func simpleUsersAPINotFoundWithCustomTransitionObjsHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errExample404

	_ = replierWithCustomTransitionObjs.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

func simpleUsersAPINotFoundHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errExample404

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

func simpleUsersAPIMultiErrorHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErrs := []error{errExampleDobValidation, errExampleNameValidation}

	_ = replier.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Errors: serverErrs,
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
		Writer:   w,
		TokenOne: mockedAccessToken,
		TokenTwo: mockedRefreshToken,
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
	serverErr := errExample404

	replierWithCustomTransitionObj.NewHTTPResponse(&reply.NewResponseRequest{
		Writer: w,
		Error:  serverErr,
	})
}

//////////////////////////////
//// Handlers Using Aides ////

func simpleUsersAPIMultiErrorUsingAideWithCustomErrorHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErrs := []error{errExampleDobValidation, errExampleNameValidation}

	_ = replierWithCustomTransitionObjError.NewHTTPMultiErrorResponse(w, serverErrs)
}

func simpleUsersAPIMultiErrorUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErrs := []error{errExampleDobValidation, errExampleNameValidation}

	_ = replier.NewHTTPMultiErrorResponse(w, serverErrs)
}

func simpleUsersAPINotFoundUsingAideHandler(w http.ResponseWriter, r *http.Request) {

	// Do something with a server
	serverErr := errExample404

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
	serverErr := errExample404

	replierWithCustomTransitionObj.NewHTTPErrorResponse(w, serverErr)
}

/////////////////////////////

func handleRequest() {
	var port string = ":8081"

	http.HandleFunc("/errors", simpleUsersAPIMultiErrorHandler)
	http.HandleFunc("/users", simpleUsersAPIHandler)
	http.HandleFunc("/users/3", simpleUsersAPINotFoundHandler)
	http.HandleFunc("/users/4", simpleUsersAPINoManifestEntryHandler)
	http.HandleFunc("/tokens/refresh", simpleTokensAPIHandler)
	http.HandleFunc("/defaults/1", simpleAPIDefaultResponseHandler)
	http.HandleFunc("/users/3/custom", simpleUsersAPINotFoundCustomReplierHandler)
	http.HandleFunc("/users/404/custom", simpleUsersAPINotFoundWithCustomTransitionObjsHandler)

	http.HandleFunc("/aides/errors", simpleUsersAPIMultiErrorUsingAideHandler)
	http.HandleFunc("/aides/errors/custom", simpleUsersAPIMultiErrorUsingAideWithCustomErrorHandler)
	http.HandleFunc("/aides/users", simpleUsersAPIUsingAideHandler)
	http.HandleFunc("/aides/users/3", simpleUsersAPINotFoundUsingAideHandler)
	http.HandleFunc("/aides/users/4", simpleUsersAPINoManifestEntryUsingAideHandler)
	http.HandleFunc("/aides/tokens/refresh", simpleTokensAPIUsingAideHandler)
	http.HandleFunc("/aides/defaults/1", simpleAPIDefaultResponseUsingAideHandler)
	http.HandleFunc("/aides/users/3/custom", simpleUsersAPINotFoundCustomReplierUsingAideHandler)

	log.Printf("Serving simple API on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	handleRequest()
}
