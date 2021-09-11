// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// TransferObjectError outlines expected methods of a transfer object error
type TransferObjectError interface {
	SetTitle(title string)
	GetTitle() string
	SetDetail(detail string)
	GetDetail() string
	SetAbout(about string)
	GetAbout() string
	SetStatusCode(status int)
	GetStatusCode() string
	SetCode(code string)
	GetCode() string
	SetMeta(meta interface{})
	GetMeta() interface{}
	RefreshTransferObject() TransferObjectError
}

// TransferObject outlines expected methods of a transfer object
type TransferObject interface {
	SetHeaders(headers map[string]string)
	SetStatusCode(code int)
	SetMeta(meta map[string]interface{})
	SetTokenOne(token string)
	SetTokenTwo(token string)
	GetWriter() http.ResponseWriter
	GetStatusCode() int
	SetWriter(writer http.ResponseWriter)
	SetErrors(transferObjectErrors []TransferObjectError)
	RefreshTransferObject() TransferObject
	SetData(data interface{})
}

const (
	// defaultResponseBody is the default response body
	defaultResponseBody = "{}"

	// defaultStatusCode is the default response status code
	defaultStatusCode = http.StatusOK

	// defaultErrorsStatusCode is the default status code for errors
	defaultErrorsStatusCode = http.StatusBadRequest
)

// Option used to build on top of default Replier features
type Option func(*Replier)

// WithTransferObject sets the base transfer object used for response
func WithTransferObject(replacementTransferObject TransferObject) Option {
	return func(r *Replier) {
		r.transferObject = replacementTransferObject
	}
}

// WithTransferObjectError sets the transfer object error used to represent
// errors in response
func WithTransferObjectError(replacementTransferObjectError TransferObjectError) Option {
	return func(r *Replier) {
		r.transferObjectError = replacementTransferObjectError
	}
}

// NewResponseRequest holds attributes for response
type NewResponseRequest struct {
	Writer     http.ResponseWriter
	Data       interface{}
	Meta       map[string]interface{}
	Headers    map[string]string
	StatusCode int
	Message    string
	Error      error
	Errors     []error
	TokenOne   string
	TokenTwo   string
}

// Replier handles managing responses
type Replier struct {
	// Error manifest used by Replier to pull corresponsing error items to build
	// response error(s).
	errorManifest ErrorManifest

	// Top-level response base with core and special attributes used to build out
	// response.
	transferObject TransferObject

	// Error object base used to shape error objects in response
	transferObjectError TransferObjectError
}

// NewReplier returns a new Replier pointer that shapes and handles both
// successful and error based responses.
//
// NOTE - Both default object(s) can be overwritten using the  chained Options, i.e.
// `WithTransferObject` or `WithTransferObjectError`
func NewReplier(manifests []ErrorManifest, options ...Option) *Replier {

	activeTransferObject := &defaultReplyTransferObject{}
	activeTransferObjectError := &defaultReplyTransferObjectError{}

	replier := Replier{
		errorManifest:       mergeManifestCollections(manifests),
		transferObject:      activeTransferObject,
		transferObjectError: activeTransferObjectError,
	}

	// Add option add-ons on replier
	for _, option := range options {
		option(&replier)
	}

	return &replier
}

// NewHTTPResponse handles generating and sending of an appropriate HTTP response body
// based response attributes.
//
// NOTE - Several assumptions have been made to simplify the process
// of response generation. The assumptions include:
//
// - An error passed in the NewResponseRequest will have a corresponding manifest entry,
// otherwise you are happy for a `500 - Internal Server Error` to be returned
//
// - Reply should only return tokens (access & refresh) with each other or by themselves
//
// - Data will be JSON encodable
//
// - The default response will be to return 200 status code if the NewResponseRequest is
// solely  passed  with a writer
func (r *Replier) NewHTTPResponse(response *NewResponseRequest) error {

	if response.Writer == nil {
		return errors.New("reply/http-response: failed to send response, no writer provided")
	}

	// Use fresh transfer object
	r.transferObject = r.transferObject.RefreshTransferObject()

	r.setUniversalAttributes(response.Writer, response.Headers, response.Meta, response.StatusCode)

	// Manage response for multi errors
	if len(response.Errors) > 0 {
		return r.generateMultiErrorResponse(response.Errors)
	}

	// Manage response for error
	if response.Error != nil {
		return r.generateErrorResponse(response.Error)
	}

	// Manage response for token
	if response.TokenOne != "" || response.TokenTwo != "" {
		return r.generateTokenResponse(response.TokenOne, response.TokenTwo)
	}

	// Manage response for data
	if response.Data != nil {
		return r.generateDataResponse(response.Data)
	}

	return r.generateDefaultResponse()
}

// generateDefaultResponse generates the default response
func (r *Replier) generateDefaultResponse() error {
	r.transferObject.SetData(defaultResponseBody)

	return sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateDataResponse generates response based on passed data
func (r *Replier) generateDataResponse(data interface{}) error {
	r.transferObject.SetData(data)

	return sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateTokenResponse generates token response on passed tokens information
func (r *Replier) generateTokenResponse(tokenOne, tokenTwo string) error {
	r.transferObject.SetTokenOne(tokenOne)
	r.transferObject.SetTokenTwo(tokenTwo)

	return sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateMultiErrorResponse generates error response for multiple
// errors
//
// NOTE - If at anytime one of the errors return a 5XX error manifest item,
// only the 5XX error will be returned
func (r *Replier) generateMultiErrorResponse(errs []error) error {

	transferObjectErrors := []TransferObjectError{}

	for _, err := range errs {
		manifestItem := r.getErrorManifestItem(err)

		if is5xx(manifestItem.StatusCode) {
			return r.sendHTTPErrorsResponse(manifestItem.StatusCode, append(
				[]TransferObjectError{},
				r.convertErrorManifestItemToTransferObjectError(manifestItem)))
		}

		transferObjectErrors = append(transferObjectErrors, r.convertErrorManifestItemToTransferObjectError(manifestItem))
	}

	statusCode := getAppropiateStatusCodeOrDefault(transferObjectErrors)

	return r.sendHTTPErrorsResponse(statusCode, transferObjectErrors)
}

// generateErrorResponse generates correct error response based on passed
// error
func (r *Replier) generateErrorResponse(err error) error {
	manifestItem := r.getErrorManifestItem(err)

	transferObjectErrors := append([]TransferObjectError{}, r.convertErrorManifestItemToTransferObjectError(manifestItem))

	return r.sendHTTPErrorsResponse(manifestItem.StatusCode, transferObjectErrors)
}

// sendHTTPErrorsResponse handles setting status code and transfer object errors before
// attempting to send response
func (r *Replier) sendHTTPErrorsResponse(statusCode int, transferObjectErrors []TransferObjectError) error {
	r.transferObject.SetStatusCode(statusCode)
	r.transferObject.SetErrors(transferObjectErrors)

	return sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// getErrorManifestItem returns the corresponding manifest Item if found,
// otherwise the internal server error is returned
func (r *Replier) getErrorManifestItem(err error) ErrorManifestItem {
	manifestItem, ok := r.errorManifest[err.Error()]
	if !ok {
		manifestItem = getInternalServertErrorManifestItem()
		log.Printf("reply/error-response: failed to find error manifest item for %v", err)
	}

	// TODO: Set default error status Code on manifest item if unset

	return manifestItem
}

// setUniversalAttributes sets the attributes that are common across all
// response types
func (r *Replier) setUniversalAttributes(writer http.ResponseWriter, headers map[string]string, meta map[string]interface{}, statusCode int) {
	r.transferObject.SetWriter(writer)
	r.setHeaders(headers)
	r.transferObject.SetMeta(meta)

	if statusCode != 0 {
		r.transferObject.SetStatusCode(statusCode)
		return
	}

	r.transferObject.SetStatusCode(defaultStatusCode)
}

// setDefaultContentType handles setting default content type to JSON if
// not already set
func (r *Replier) setDefaultContentType() {
	if r.transferObject.GetWriter().Header().Get("Content-type") == "" {
		r.transferObject.GetWriter().Header().Set("Content-type", "application/json")
	}
}

// setHeaders handles setting headers on the writer. Existing headers should not
// be affected unless they share the header key
func (r *Replier) setHeaders(h map[string]string) {

	r.setDefaultContentType()

	if h == nil {
		return
	}

	for headerKey, headerValue := range h {
		r.transferObject.GetWriter().Header().Set(headerKey, headerValue)
	}
}

// convertErrorManifestItemToTransferObjectError converts manifest error item to valid
// transfer object error
func (r *Replier) convertErrorManifestItemToTransferObjectError(errorItem ErrorManifestItem) TransferObjectError {

	// Use fresh transfer object error
	convertedError := r.transferObjectError.RefreshTransferObject()

	convertedError.SetTitle(errorItem.Title)
	convertedError.SetDetail(errorItem.Detail)
	convertedError.SetAbout(errorItem.About)
	convertedError.SetCode(errorItem.Code)
	convertedError.SetStatusCode(errorItem.StatusCode)
	convertedError.SetMeta(errorItem.Meta)

	return convertedError
}

// getAppropiateStatusCodeOrDefault loops through collection of transfer object errors (first to last), and
// attempts to pull and convert status code (string).
//
// NOTE - If error occurs the next element will be attempted. In the event no elements are left, the default
// error status code (400) will be returned
func getAppropiateStatusCodeOrDefault(transferObjectErrors []TransferObjectError) int {

	for _, transferObjectError := range transferObjectErrors {

		statusCode, err := strconv.Atoi(transferObjectError.GetStatusCode())
		if err != nil {
			continue
		}

		return statusCode
	}

	return defaultErrorsStatusCode
}

// is5xx returns whether status code is a 5xx
func is5xx(statusCode int) bool {
	if statusCode >= 500 && statusCode <= 599 {
		return true
	}

	return false
}

// sendHTTPResponse handles sending response based on the transfer object
func sendHTTPResponse(writer http.ResponseWriter, transferObject TransferObject) error {

	writer.WriteHeader(transferObject.GetStatusCode())
	err := json.NewEncoder(writer).Encode(transferObject)
	if err != nil {
		return fmt.Errorf("reply/http-response: failed to encode transfer object with %v", err)
	}

	return nil
}

// mergeManifestCollections handles merging the passed manifests into a singular
// manifest
func mergeManifestCollections(manifests []ErrorManifest) ErrorManifest {

	mergedManifests := make(ErrorManifest)

	for _, manifest := range manifests {
		getManifestItems(manifest, mergedManifests)
	}

	return mergedManifests
}

// getManifestItems pulls the key and items from the manifest and inserts into final manifest
func getManifestItems(manifest ErrorManifest, finalManifest ErrorManifest) {

	for key, item := range manifest {
		finalManifest[key] = item
	}
}

// getInternalServertErrorManifestItem returns typical 500 error with text and message
func getInternalServertErrorManifestItem() ErrorManifestItem {
	return ErrorManifestItem{Title: "Internal Server Error", StatusCode: http.StatusInternalServerError}
}

/////////////////////////////////////////////////
//////////////// Response Aides /////////////////
// Response aides simplify how users interact
// with this library to create their success and
// error driven responses.

// ResponseAttributes used to add additional attributes to
// response aides
type ResponseAttributes func(*NewResponseRequest)

// WithHeaders adds passed headers on to the generated response
func WithHeaders(headers map[string]string) ResponseAttributes {
	return func(r *NewResponseRequest) {
		r.Headers = headers
	}
}

// WithMeta adds passed meta data on to the generated response
func WithMeta(meta map[string]interface{}) ResponseAttributes {
	return func(r *NewResponseRequest) {
		r.Meta = meta
	}
}

// NewHTTPMultiErrorResponse this response aide is used to create
// a multi error response. It will utilise the manifest
// declared when creating its base replier to pull all corresponding
// error manifest items.
//
// With this aide, if desired, you can add additional attributes by using the
// WithHeaders and/ or WithMeta optional response attributes.
//
// NOTE - If ANY of the passed errors do not have a manifest entry, a single
// 500 error will be returned.
func (r *Replier) NewHTTPMultiErrorResponse(w http.ResponseWriter, errs []error, attributes ...ResponseAttributes) error {

	request := NewResponseRequest{
		Writer: w,
		Errors: errs,
	}

	// Add attributes to response request
	for _, attribute := range attributes {
		attribute(&request)
	}

	return r.NewHTTPResponse(&request)
}

// NewHTTPErrorResponse this response aide is used to create
// response explicitly for errors. It will utilise the manifest
// declared when creating its base replier.
//
// With this aide, if desired, you can add additional attributes by using the
// WithHeaders and/ or WithMeta optional response attributes.
//
// NOTE - If the passed error doesn't have a manifest entry, a 500 error will
// be returned.
func (r *Replier) NewHTTPErrorResponse(w http.ResponseWriter, err error, attributes ...ResponseAttributes) error {

	request := NewResponseRequest{
		Writer: w,
		Error:  err,
	}

	// Add attributes to response request
	for _, attribute := range attributes {
		attribute(&request)
	}

	return r.NewHTTPResponse(&request)
}

// NewHTTPDataResponse this response aide is used to create
// "successful" response. "Successful" response are responses
// that contain some sort of data that will be returned to the
// consumer.
//
// With this aide, if desired, you can add additional attributes by using the
// WithHeaders and/ or WithMeta optional response attributes.
func (r *Replier) NewHTTPDataResponse(w http.ResponseWriter, statusCode int, data interface{}, attributes ...ResponseAttributes) error {

	request := NewResponseRequest{
		Writer:     w,
		Data:       data,
		StatusCode: statusCode,
	}

	// Add attributes to response request
	for _, attribute := range attributes {
		attribute(&request)
	}

	return r.NewHTTPResponse(&request)
}

// NewHTTPBlankResponse this response aide is used to create
// a "blank" response. A blank response is one that contains
// the response body "{}".
//
// With this aide, if desired, you can add additional attributes by using the
// WithHeaders and/ or WithMeta optional response attributes.
func (r *Replier) NewHTTPBlankResponse(w http.ResponseWriter, statusCode int, attributes ...ResponseAttributes) error {

	request := NewResponseRequest{
		Writer:     w,
		StatusCode: statusCode,
	}

	// Add attributes to response request
	for _, attribute := range attributes {
		attribute(&request)
	}

	return r.NewHTTPResponse(&request)
}

// NewHTTPTokenResponse this response aide is used to create
// the response for token(s). If the desired behaviour is to return
// a single token, pass an empty string in the token not to be
// included in the response.
//
// With this aide, if desired, you can add additional attributes by using the
// WithHeaders and/ or WithMeta optional response attributes.
//
// NOTE - At least one of the tokens must be specified or an error will
// be returned
func (r *Replier) NewHTTPTokenResponse(w http.ResponseWriter, statusCode int, tokenOne, tokenTwo string, attributes ...ResponseAttributes) error {

	if isEmpty(tokenOne) && isEmpty(tokenTwo) {
		return errors.New("reply/http-token-aide: failed at least one token must be returned")
	}

	request := NewResponseRequest{
		Writer:     w,
		StatusCode: statusCode,
		TokenOne:   tokenOne,
		TokenTwo:   tokenTwo,
	}

	// Add attributes to response request
	for _, attribute := range attributes {
		attribute(&request)
	}

	return r.NewHTTPResponse(&request)
}

// isEmpty checks if the passed string is empty
func isEmpty(s string) bool {
	return s == ""
}
