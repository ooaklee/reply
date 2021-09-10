// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply

import (
	"net/http"
	"strconv"
)

// ErrorManifestItem holds the message and status code for an error response
type ErrorManifestItem struct {

	// Title holds the text summary returned in the response's error response
	//
	// NOTE:
	//
	// - The most effective title are short, sweet and easy to consume
	//
	// - This message will be seen by the consuming client, be mindful of
	// the amount of information you divulge
	Title string

	// Detail holds a more descriptive brief returned in the response's error response
	//
	// NOTE:
	//
	// - Detail will give a deeper level of context, while being mindful of length
	//
	// - Like the title message will be seen by the consuming client, be mindful of
	// the amount of information you divulge
	Detail string

	// StatusCode holds the HTTP status code that best relates to the response.
	// For more information on status codes, https://httpstatuses.com/.
	StatusCode int

	// About holds the a URL that gives further insight into the error
	About string

	// Code holds the internal application error code, if appicable, thst is used to
	// help debuggers better identify error
	Code string

	// Meta contains additional meta-information about the that can be shared to
	// consumer
	Meta interface{}
}

// ErrorManifest holds error reference (string) with its corresponding
// manifest item (message & status code) which it returned in the response
type ErrorManifest map[string]ErrorManifestItem

// Error holds attributes often used to give additional
// context when unexpected behaviour occurs
type Error struct {

	// Title a short summary of the problem
	Title string `json:"title,omitempty"`

	// Detail a description of the error
	Detail string `json:"detail,omitempty"`

	// About holds the link that gives further insight into the error
	About string `json:"about,omitempty"`

	// Status the HTTP status associated with error
	Status string `json:"status,omitempty"`

	// Code internal error code used to reference error
	Code string `json:"code,omitempty"`

	// Meta contains additional meta-information about the error
	Meta interface{} `json:"meta,omitempty"`
}

// SetTitle adds title to error
func (e *Error) SetTitle(title string) {
	e.Title = title
}

// GetTitle returns error's title
func (e *Error) GetTitle() string {
	return e.Title
}

// SetDetail adds detail to error
func (e *Error) SetDetail(detail string) {
	e.Detail = detail
}

// GetDetail return error's detail
func (e *Error) GetDetail() string {
	return e.Detail
}

// SetAbout adds about to error
func (e *Error) SetAbout(about string) {
	e.About = about
}

// GetAbout return error's about
func (e *Error) GetAbout() string {
	return e.About
}

// SetStatusCode converts and add http status code to error
func (e *Error) SetStatusCode(status int) {
	e.Status = strconv.Itoa(status)
}

// GetStatusCode returns error's HTTP status code
func (e *Error) GetStatusCode() string {
	return e.Status
}

// SetCode adds internal code to error
func (e *Error) SetCode(code string) {
	e.Code = code
}

// GetCode returns error's internal code
func (e *Error) GetCode() string {
	return e.Code
}

// SetMeta adds meta property to error
func (e *Error) SetMeta(meta interface{}) {
	e.Meta = meta
}

// GetMeta returns error's meta property
func (e *Error) GetMeta() interface{} {
	return e.Meta
}

// RefreshTransferObject returns an empty instance of transfer object
// error
func (e *Error) RefreshTransferObject() TransferObjectError {
	return &Error{}
}

// defaultReplyTransferObject handles structing response for client
// consumption
type defaultReplyTransferObject struct {
	HTTPWriter   http.ResponseWriter    `json:"-"`
	Headers      map[string]string      `json:"-"`
	StatusCode   int                    `json:"-"`
	Errors       []TransferObjectError  `json:"errors,omitempty"`
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

// SetErrors assigns the passed transfer object errors to the transfer object
func (t *defaultReplyTransferObject) SetErrors(transferObjectErrors []TransferObjectError) {
	t.Errors = transferObjectErrors
}
