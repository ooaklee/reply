package reply

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// TransferObject outlines expected methods of a transfer object
type TransferObject interface {
	SetHeaders(headers map[string]string)
	SetStatusCode(code int)
	SetMeta(meta map[string]interface{})
	SetAccessToken(token string)
	SetRefreshToken(token string)
	GetWriter() http.ResponseWriter
	GetStatusCode() int
	SetWriter(writer http.ResponseWriter)
	SetStatus(transferObjectStatus *TransferObjectStatus)
	RefreshTransferObject() TransferObject
	SetData(data interface{})
}

const (
	// defaultResponseBody returns default response body
	defaultResponseBody = "{}"

	// defaultStatusCode returns default response status code
	defaultStatusCode = http.StatusOK
)

// Option used to build ontop of default features
type Option func(*Replier)

// WithTransferObject overwrites the transfer object used for response
func WithTransferObject(replacementTransferObject TransferObject) Option {
	return func(r *Replier) {
		r.transferObject = replacementTransferObject
	}
}

// NewResponseRequest holds attributes for response
type NewResponseRequest struct {
	Writer       http.ResponseWriter
	Data         interface{}
	Meta         map[string]interface{}
	Headers      map[string]string
	StatusCode   int
	Message      string
	Error        error
	AccessToken  string
	RefreshToken string
}

// Replier handles managing responses
type Replier struct {
	errorManifest  ErrorManifest
	transferObject TransferObject
}

// NewReplier creates a replier
func NewReplier(manifests []ErrorManifest, options ...Option) *Replier {

	activeTransferObject := &defaultReplyTransferObject{}

	replier := Replier{
		errorManifest:  mergeManifestCollections(manifests...),
		transferObject: activeTransferObject,
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

	// Manage response for error
	if response.Error != nil {
		r.generateErrorResponse(response.Error)
		return nil
	}

	// Manage response for token
	if response.AccessToken != "" || response.RefreshToken != "" {
		r.generateTokenResponse(response.AccessToken, response.RefreshToken, response.StatusCode)
		return nil
	}

	// Manage response for data
	if response.Data != nil {
		r.generateDataResponse(response.Data, response.StatusCode)
		return nil
	}

	r.generateDefaultResponse()
	return nil
}

// generateDefaultResponse generates the default response
func (r *Replier) generateDefaultResponse() {
	r.transferObject.SetStatusCode(defaultStatusCode)
	r.transferObject.SetData(defaultResponseBody)

	sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateDataResponse generates response based on passed data
func (r *Replier) generateDataResponse(data interface{}, statusCode int) {
	r.transferObject.SetData(data)

	if statusCode == 0 {
		r.transferObject.SetStatusCode(defaultStatusCode)
	}

	sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateTokenResponse generates token response on passed tokens information
func (r *Replier) generateTokenResponse(accessToken, refreshToken string, statusCode int) {
	r.transferObject.SetAccessToken(accessToken)
	r.transferObject.SetRefreshToken(refreshToken)

	if statusCode == 0 {
		r.transferObject.SetStatusCode(defaultStatusCode)
	}

	sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// generateErrorResponse generates correct error response based on passed
// error
func (r *Replier) generateErrorResponse(err error) {
	manifestItem, ok := r.errorManifest[err.Error()]
	if !ok {
		manifestItem = getInternalServertErrorManifestItem()
	}

	transferObjectStatus := &TransferObjectStatus{}
	transferObjectStatus.SetMessage(manifestItem.Message)

	// Overwrite status code
	r.transferObject.SetStatusCode(manifestItem.StatusCode)
	r.transferObject.SetStatus(transferObjectStatus)

	sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
}

// setUniversalAttributes sets the attributes that are common across all
// response types
func (r *Replier) setUniversalAttributes(writer http.ResponseWriter, headers map[string]string, meta map[string]interface{}, statusCode int) {
	r.transferObject.SetWriter(writer)
	r.setHeaders(headers)
	r.transferObject.SetMeta(meta)
	r.transferObject.SetStatusCode(statusCode)
}

// setDefaultContentType handles setting default content type to JSON if
// not already set
func (r *Replier) setDefaultContentType() {
	if r.transferObject.GetWriter().Header().Get("Content-type") == "" {
		r.transferObject.GetWriter().Header().Set("Content-type", "application/json")
	}
}

// setHeaders handles setting headers on writer. Existing headers hould not
// be affeted unless they share the header key
func (r *Replier) setHeaders(h map[string]string) {

	r.setDefaultContentType()

	if h == nil {
		return
	}

	for headerKey, headerValue := range h {
		r.transferObject.GetWriter().Header().Set(headerKey, headerValue)
	}
}

// sendHTTPResponse handles sending response based on the transfer object
func sendHTTPResponse(writer http.ResponseWriter, transferObject TransferObject) {

	writer.WriteHeader(transferObject.GetStatusCode())
	err := json.NewEncoder(writer).Encode(transferObject)
	if err == nil {
		return
	}

	log.Printf("reply/http-response: failed to encode transfer object with %v", err)

	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte{})
}

// mergeManifestCollections handles merges the passed manifests into a singular
// map
func mergeManifestCollections(manifests ...ErrorManifest) ErrorManifest {

	mergedManifests := make(ErrorManifest)

	for _, manifest := range manifests {
		key, value := getManifestAttributes(manifest)
		mergedManifests[key] = *value
	}

	return mergedManifests
}

// getManifestAttributes returns key and value for pass manifest
func getManifestAttributes(manifest ErrorManifest) (key string, value *ErrorManifestItem) {

	for k, v := range manifest {
		key = k
		value = &v
	}

	return key, value
}

// getInternalServertErrorManifestItem returns typical 500 error with text and message
func getInternalServertErrorManifestItem() ErrorManifestItem {
	return ErrorManifestItem{Message: "Internal Server Error", StatusCode: http.StatusInternalServerError}
}
