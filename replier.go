package reply

import (
	"encoding/json"
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

// NewHTTPResponse handles generating and sending of an appropiate HTTP response body
// based response attributes.
//
// NOTE - A number of a assumptions have been made to simplify the process
// of response generation. The assumptions include:
//
// - Responses with a StatusCode `NOT` between 200 - 299, and 300-301 will be
// deemed as an error response.
func (r *Replier) NewHTTPResponse(response *NewResponseRequest) error {

	r.transferObject.SetHeaders(response.Headers)
	r.transferObject.SetMeta(response.Meta)
	r.transferObject.SetWriter(response.Writer)
	r.transferObject.SetStatusCode(response.StatusCode)

	// Manage response for error
	if response.Error != nil {

		manifestItem, ok := r.errorManifest[response.Error.Error()]
		if !ok {
			manifestItem = getInternalServertErrorManifestItem()
		}

		transferObjectStatus := &TransferObjectStatus{}
		transferObjectStatus.SetMessage(manifestItem.Message)

		// Overwrite status code
		r.transferObject.SetStatusCode(manifestItem.StatusCode)
		r.transferObject.SetStatus(transferObjectStatus)

		sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
		r.transferObject = r.transferObject.RefreshTransferObject()
		return nil
	}

	// Manage response for token
	if response.AccessToken != "" || response.RefreshToken != "" {
		r.transferObject.SetAccessToken(response.AccessToken)
		r.transferObject.SetRefreshToken(response.RefreshToken)

		sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
		r.transferObject = r.transferObject.RefreshTransferObject()
		return nil
	}

	// Manage response for data
	if response.Data != nil {
		r.transferObject.SetData(response.Data)

		if response.StatusCode == 0 {
			r.transferObject.SetStatusCode(defaultStatusCode)
		}

		sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
		r.transferObject = r.transferObject.RefreshTransferObject()
		return nil
	}

	// Set Default response
	r.transferObject.SetStatusCode(defaultStatusCode)
	r.transferObject.SetData(defaultResponseBody)

	sendHTTPResponse(r.transferObject.GetWriter(), r.transferObject)
	r.transferObject = r.transferObject.RefreshTransferObject()
	return nil
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
