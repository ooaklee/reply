// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

// WithTransferObject supplies a response prototype. RefreshTransferObject must
// return independent state and be safe to call concurrently without mutating
// this prototype when the Replier is shared between requests.
func WithTransferObject(replacementTransferObject TransferObject) Option {
	return func(r *Replier) {
		r.transferObject = replacementTransferObject
	}
}

// WithTransferObjectError supplies an error prototype. Its refresh method must
// return independent objects and be concurrency-safe for a shared Replier.
func WithTransferObjectError(replacementTransferObjectError TransferObjectError) Option {
	return func(r *Replier) {
		r.transferObjectError = replacementTransferObjectError
	}
}

// NewResponseRequest holds attributes for response
type NewResponseRequest struct {
	// Context is optional request context for safe failure observation; it does
	// not alter response fields, status, or the writer's cancellation semantics.
	Context context.Context
	// Writer is the caller-owned destination, required for every response.
	Writer http.ResponseWriter
	// Data is the success payload; nil retains the existing blank-response behavior.
	Data interface{}
	// Meta contains caller-approved public response metadata.
	Meta map[string]interface{}
	// Headers contains caller-approved HTTP headers applied to this response.
	Headers map[string]string
	// StatusCode selects the success status; mapped errors use their manifest status.
	StatusCode int
	// Message is retained for source compatibility; built-in responses do not use it.
	Message string
	// Error supplies one error tree when Errors is empty.
	Error error
	// Errors supplies error trees with precedence over Error and success payloads.
	Errors []error
	// TokenOne is the optional access token returned by a successful token response.
	TokenOne string
	// TokenTwo is the optional refresh token returned by a successful token response.
	TokenTwo string
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

	// unmappedErrorObserver reports incomplete resolution without raw diagnostics.
	unmappedErrorObserver func(context.Context)
}

// WithUnmappedErrorObserver observes one unresolved response using its request
// context, or context.Background when none was supplied. No error text or values
// are passed and default behavior is silent. The callback runs synchronously,
// must be concurrency-safe and should not panic or block response delivery.
// Use it for request-correlated metrics; redact any context-derived log fields.
func WithUnmappedErrorObserver(observer func(context.Context)) Option {
	return func(r *Replier) { r.unmappedErrorObserver = observer }
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
		if option != nil {
			option(&replier)
		}
	}

	return &replier
}

// NewHTTPResponse handles generating and sending of an appropriate HTTP response body
// based response attributes.
// All response state is request-local; the replier can be shared if its custom
// prototypes and observer follow their documented concurrency contracts.
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
	if response == nil {
		return errors.New("reply/http-response: no response request provided")
	}
	if isNilValue(response.Writer) {
		return errors.New("reply/http-response: failed to send response, no writer provided")
	}
	if r == nil || isNilValue(r.transferObject) || isNilValue(r.transferObjectError) {
		return errors.New("reply/http-response: invalid response prototype")
	}
	// Keep the reusable manifest and prototypes immutable. All writer and body
	// state belongs to this request, including for concurrent calls.
	local := *r
	local.transferObject = r.transferObject.RefreshTransferObject()
	if isNilValue(local.transferObject) || samePrototype(local.transferObject, r.transferObject) {
		return errors.New("reply/http-response: response prototype must return independent non-nil state")
	}
	return local.writeHTTPResponse(response)
}

// writeHTTPResponse applies existing precedence to request-local response state.
func (r *Replier) writeHTTPResponse(response *NewResponseRequest) error {

	r.setUniversalAttributes(response.Writer, response.Headers, response.Meta, response.StatusCode)

	// Manage response for multi errors
	if len(response.Errors) > 0 {
		return r.generateMultiErrorResponse(response.Context, response.Errors)
	}

	// Manage response for error
	if response.Error != nil {
		return r.generateMultiErrorResponse(response.Context, []error{response.Error})
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

// generateMultiErrorResponse resolves every independent failure before emitting
// manifest-driven errors. A missing classification invalidates the collection;
// diagnostics are never passed to a logger or to the transfer object.
func (r *Replier) generateMultiErrorResponse(ctx context.Context, errs []error) error {
	keys, known := r.ResolveErrors(errs...)
	if !known {
		if r.unmappedErrorObserver != nil {
			if ctx == nil {
				ctx = context.Background()
			}
			r.unmappedErrorObserver(ctx)
		}
		return r.sendManifestErrors([]ErrorManifestItem{getInternalServertErrorManifestItem()})
	}
	items := make([]ErrorManifestItem, 0, len(keys))
	for _, key := range keys {
		item := r.errorManifest[key]
		setDefaultStatusCode(&item)
		if is5xx(item.StatusCode) {
			return r.sendManifestErrors([]ErrorManifestItem{item})
		}
		items = append(items, item)
	}
	return r.sendManifestErrors(items)
}

// sendManifestErrors encodes approved items using per-error prototype instances.
// All factories complete before status/body writes, so an invalid factory does
// not commit a partial response. The first item's status retains precedence.
func (r *Replier) sendManifestErrors(items []ErrorManifestItem) error {
	converted := make([]TransferObjectError, 0, len(items))
	for _, item := range items {
		value, err := r.convertErrorManifestItemToTransferObjectError(item)
		if err != nil {
			return err
		}
		converted = append(converted, value)
	}
	// The manifest, not a custom serializer's status getter, owns HTTP status.
	return r.sendHTTPErrorsResponse(items[0].StatusCode, converted)
}

// sendHTTPErrorsResponse handles setting status code and transfer object errors before
// attempting to send response
func (r *Replier) sendHTTPErrorsResponse(statusCode int, transferObjectErrors []TransferObjectError) error {
	r.transferObject.SetStatusCode(statusCode)
	r.transferObject.SetErrors(transferObjectErrors)

	return sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// setDefaultStatusCode sets the error manifest item's status code to default error
// code value if it is not already set (non-zero)
func setDefaultStatusCode(item *ErrorManifestItem) {
	if item.StatusCode == 0 {
		item.StatusCode = defaultErrorsStatusCode
	}
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
func (r *Replier) convertErrorManifestItemToTransferObjectError(errorItem ErrorManifestItem) (TransferObjectError, error) {

	// Use fresh transfer object error
	convertedError := r.transferObjectError.RefreshTransferObject()
	if isNilValue(convertedError) || samePrototype(convertedError, r.transferObjectError) {
		return nil, errors.New("reply/http-response: error prototype must return independent non-nil state")
	}

	convertedError.SetTitle(errorItem.Title)
	convertedError.SetDetail(errorItem.Detail)
	convertedError.SetAbout(errorItem.About)
	convertedError.SetCode(errorItem.Code)
	convertedError.SetStatusCode(errorItem.StatusCode)
	convertedError.SetMeta(errorItem.Meta)

	return convertedError, nil
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
		return fmt.Errorf("reply/http-response: failed to encode transfer object with %w", err)
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

// getManifestItems pulls the error and items from the manifest and inserts into final manifest
func getManifestItems(manifest ErrorManifest, finalManifest ErrorManifest) {

	for err, item := range manifest {
		finalManifest[err] = item
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

// WithContext supplies request context for an unmapped-error observer. Context
// values are not serialized or logged by reply and do not establish authority.
func WithContext(ctx context.Context) ResponseAttributes {
	return func(r *NewResponseRequest) { r.Context = ctx }
}

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
		if attribute != nil {
			attribute(&request)
		}
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
		if attribute != nil {
			attribute(&request)
		}
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
		if attribute != nil {
			attribute(&request)
		}
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
		if attribute != nil {
			attribute(&request)
		}
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
		if attribute != nil {
			attribute(&request)
		}
	}

	return r.NewHTTPResponse(&request)
}

// isEmpty checks if the passed string is empty
func isEmpty(s string) bool {
	return s == ""
}
