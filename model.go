// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply

import "net/http"

// ErrorManifestItem holds the message and status code for a error response
type ErrorManifestItem struct {

	// Message holds the text returned in the response's status.message path
	//
	// NOTE:
	//
	// - The most effective messages are short, sweet and easy to consume
	//
	// - This message will be seeing by the consuming client, be mindful of
	// how much information you divulge
	Message string

	// StatusCode holds the HTTP status code that best relates to the response.
	// For more information on status codes, https://httpstatuses.com/.
	StatusCode int
}

// ErrorManifest holds error reference (string) with its corresponding
// manifest item (message & status code) which it returned in the response
type ErrorManifest map[string]ErrorManifestItem

// TransferObjectStatus holds attributes often used to give additional
// context in responses
type TransferObjectStatus struct {
	Errors  []Error `json:"errors,omitempty"`
	Message string  `json:"message,omitempty"`
}

// SetMessage adds message to transfer object status
func (s *TransferObjectStatus) SetMessage(message string) {
	s.Message = message
}

// Error holds the associated code and detail for errors passed
type Error struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}

// defaultReplyTransferObject handles structing response for client
// consumption
type defaultReplyTransferObject struct {
	HTTPWriter   http.ResponseWriter    `json:"-"`
	Headers      map[string]string      `json:"-"`
	StatusCode   int                    `json:"-"`
	Status       *TransferObjectStatus  `json:"status,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Data         interface{}            `json:"data,omitempty"`
	AccessToken  string                 `json:"access_token,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty"`
}

// SetHeaders adds headers to transfer object
// TODO: Think about any validation that can be added
func (t *defaultReplyTransferObject) SetHeaders(headers map[string]string) {
	t.Headers = headers
}

// SetHeaders adds status code to transfer object
// TODO: Think about any validation that can be added
func (t *defaultReplyTransferObject) SetStatusCode(code int) {
	t.StatusCode = code
}

// SetMeta adds meta property to transfer object
func (t *defaultReplyTransferObject) SetMeta(meta map[string]interface{}) {
	t.Meta = meta
}

// SetWriter adds writer to transfer object
func (t *defaultReplyTransferObject) SetWriter(writer http.ResponseWriter) {
	t.HTTPWriter = writer
}

// SetAccessToken adds token to access token property on transfer object
func (t *defaultReplyTransferObject) SetAccessToken(token string) {
	t.AccessToken = token
}

// SetRefreshToken adds token to refresh token property on transfer object
func (t *defaultReplyTransferObject) SetRefreshToken(token string) {
	t.RefreshToken = token
}

// GetWriter returns the writer assigned with the transfer object
func (t *defaultReplyTransferObject) GetWriter() http.ResponseWriter {
	return t.HTTPWriter
}

// GetStatusCode returns the status code assigned to the transfer object
func (t *defaultReplyTransferObject) GetStatusCode() int {
	return t.StatusCode
}

// SetData adds passed data to the transfer object
func (t *defaultReplyTransferObject) SetData(data interface{}) {
	t.Data = data
}

// RefreshTransferObject returns an empty instance of transfer object
func (t *defaultReplyTransferObject) RefreshTransferObject() TransferObject {
	return &defaultReplyTransferObject{}
}

// SetStatus assigns the passed transfer object status to the transfer object
func (t *defaultReplyTransferObject) SetStatus(transferObjectStatus *TransferObjectStatus) {
	t.Status = transferObjectStatus
}
