// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ooaklee/reply"
	"github.com/stretchr/testify/assert"
)

// user mock data object to return
type user struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func TestReplier_NewHTTPResponse(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		request            reply.NewResponseRequest
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		///////////////////////////////
		/////// Blank Response ////////
		{
			name:               "Success - Blank response (default)",
			manifests:          getEmptyErrorManifest(),
			request:            reply.NewResponseRequest{},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with different status code",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Headers:    getReplyFormattedHeader(),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Meta:       getReplyFormattedMeta(),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Headers:    getReplyFormattedHeader(),
				Meta:       getReplyFormattedMeta(),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		//////// Data Response ////////
		{
			name:      "Success - Data response",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Data:       getTestUser(),
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Data:       getTestUser(),
				StatusCode: 201,
				Headers:    getReplyFormattedHeader(),
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 201,
				Data:       getTestUser(),
				Meta:       getReplyFormattedMeta(),
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 201,
				Headers:    getReplyFormattedHeader(),
				Meta:       getReplyFormattedMeta(),
				Data:       getTestUser(),
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		/////// Token Response ////////
		{
			name:      "Success - Token response",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
				StatusCode: 200,
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Token response (single token)",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				TokenOne:   "test-token-1",
				StatusCode: 200,
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getSingleTokenResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Token response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
				StatusCode: 200,
				Headers:    getReplyFormattedHeader(),
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Token response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 200,
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
				Meta:       getReplyFormattedMeta(),
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Token response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 200,
				Headers:    getReplyFormattedHeader(),
				Meta:       getReplyFormattedMeta(),
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		/////// Error Response ////////
		{
			name:      "Failure - Error response",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Error: getExampleErrorOne(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error: getExampleErrorOne(),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOne()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Multi Error response",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors: getMultiErrors(),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Failure -  Multi Error response missing error",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Errors: getMultiErrorsWithMissingErr(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Additional Headers",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error:   getExampleErrorOne(),
				Headers: getReplyFormattedHeader(),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOne()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Additional Headers",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors:  getMultiErrors(),
				Headers: getReplyFormattedHeader(),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Meta-Information",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error: getExampleErrorOne(),
				Meta:  getReplyFormattedMeta(),
			},
			expectedStatusCode: 404,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Meta-Information",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors: getMultiErrors(),
				Meta:   getReplyFormattedMeta(),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Meta-Information & Additional Header",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error:   getExampleErrorOne(),
				Headers: getReplyFormattedHeader(),
				Meta:    getReplyFormattedMeta(),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Meta-Information & Additional Header",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors:  getMultiErrors(),
				Headers: getReplyFormattedHeader(),
				Meta:    getReplyFormattedMeta(),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		////// Response Ranking ///////
		{
			name:      "Success - Multi Error response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors:     getMultiErrors(),
				Error:      getExampleErrorOne(),
				Data:       getTestUser(),
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error:      getExampleErrorOne(),
				Data:       getTestUser(),
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOne()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Data:       getTestUser(),
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Default response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			test.request.Writer = w

			replier.NewHTTPResponse(&test.request)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPErrorResponseAide(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		passedError        error
		responseAttributes []reply.ResponseAttributes
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Failure - Error response",
			manifests:          getEmptyErrorManifest(),
			passedError:        getExampleErrorOne(),
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOne()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:        "Success - Error response with Additional Headers",
			manifests:   getDefaultErrorManifest(),
			passedError: getExampleErrorOne(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOne()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:        "Success - Error response with Meta-Information",
			manifests:   getDefaultErrorManifest(),
			passedError: getExampleErrorOne(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 404,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:        "Success - Error response with Meta-Information & Additional Header",
			manifests:   getDefaultErrorManifest(),
			passedError: getExampleErrorOne(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			if len(test.responseAttributes) > 0 {
				replier.NewHTTPErrorResponse(w, test.passedError, test.responseAttributes...)
			} else {
				replier.NewHTTPErrorResponse(w, test.passedError)
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPMultiErrorResponseAide(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		passedErrors       []error
		responseAttributes []reply.ResponseAttributes
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Success - Multi Error response",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Failure -  Multi Error response missing error",
			manifests:          getEmptyErrorManifest(),
			passedErrors:       getMultiErrorsWithMissingErr(),
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},

		{
			name:         "Success - Multi error response with Additional Headers",
			manifests:    getDefaultErrorManifest(),
			passedErrors: getMultiErrors(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},

		{
			name:         "Success - Multi error response with Meta-Information",
			manifests:    getDefaultErrorManifest(),
			passedErrors: getMultiErrors(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},

		{
			name:         "Success - Multi error response with Meta-Information & Additional Header",
			manifests:    getDefaultErrorManifest(),
			passedErrors: getMultiErrors(),
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			if len(test.responseAttributes) > 0 {
				replier.NewHTTPMultiErrorResponse(w, test.passedErrors, test.responseAttributes...)
			} else {
				replier.NewHTTPMultiErrorResponse(w, test.passedErrors)
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPDataResponseAide(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		passedStatusCode   int
		passedData         interface{}
		responseAttributes []reply.ResponseAttributes
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:             "Success - Data response",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 201,
			passedData:       getTestUser(),

			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Data response with Additional Headers",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 201,
			passedData:       getTestUser(),

			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:             "Success - Data response with Meta-Information",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 201,
			passedData:       getTestUser(),

			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Data response with Meta-Information & Additional Header",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 201,
			passedData:       getTestUser(),

			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getDataResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			if len(test.responseAttributes) > 0 {
				replier.NewHTTPDataResponse(w, test.passedStatusCode, test.passedData, test.responseAttributes...)
			} else {
				replier.NewHTTPDataResponse(w, test.passedStatusCode, test.passedData)
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPTokenResponseAide(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		passedStatusCode   int
		passedTokenOne     string
		passedTokenTwo     string
		responseAttributes []reply.ResponseAttributes
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Success - Token response",
			manifests:          getEmptyErrorManifest(),
			passedTokenOne:     "test-token-1",
			passedTokenTwo:     "test-token-2",
			passedStatusCode:   200,
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Token response (single token)",
			manifests:          getEmptyErrorManifest(),
			passedTokenOne:     "test-token-1",
			passedStatusCode:   200,
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getSingleTokenResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Token response with Additional Headers",
			manifests:        getEmptyErrorManifest(),
			passedTokenOne:   "test-token-1",
			passedTokenTwo:   "test-token-2",
			passedStatusCode: 200,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:             "Success - Token response with Meta-Information",
			manifests:        getEmptyErrorManifest(),
			passedTokenOne:   "test-token-1",
			passedTokenTwo:   "test-token-2",
			passedStatusCode: 200,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Token response with Meta-Information & Additional Header",
			manifests:        getEmptyErrorManifest(),
			passedTokenOne:   "test-token-1",
			passedTokenTwo:   "test-token-2",
			passedStatusCode: 200,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getFullTokenResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			if len(test.responseAttributes) > 0 {
				replier.NewHTTPTokenResponse(w, test.passedStatusCode, test.passedTokenOne, test.passedTokenTwo, test.responseAttributes...)
			} else {
				replier.NewHTTPTokenResponse(w, test.passedStatusCode, test.passedTokenOne, test.passedTokenTwo)
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPBlankResponseAide(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		passedStatusCode   int
		responseAttributes []reply.ResponseAttributes
		transferObject     reply.TransferObject
		transferObjecError reply.TransferObjectError
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Success - Blank response (default)",
			manifests:          getEmptyErrorManifest(),
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Blank response with different status code",
			manifests:          getEmptyErrorManifest(),
			passedStatusCode:   302,
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Blank response with Additional Headers",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 302,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:             "Success - Blank response with Meta-Information",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 302,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:             "Success - Blank response with Meta-Information & Additional Header",
			manifests:        getEmptyErrorManifest(),
			passedStatusCode: 302,
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getBlankResponseWithMetaBody()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			var replier *reply.Replier

			switch {
			case test.transferObject != nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject), reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject != nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObject(test.transferObject))
			case test.transferObject == nil && test.transferObjecError != nil:
				replier = reply.NewReplier(test.manifests, reply.WithTransferObjectError(test.transferObjecError))
			case test.transferObject == nil && test.transferObjecError == nil:
				replier = reply.NewReplier(test.manifests)
			}

			if len(test.responseAttributes) > 0 {
				replier.NewHTTPBlankResponse(w, test.passedStatusCode, test.responseAttributes...)
			} else {
				replier.NewHTTPBlankResponse(w, test.passedStatusCode)
			}

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

// stringWithNewLine appends new line to passed string
func stringWithNewLine(s string) string {
	return fmt.Sprintf("%s\n", s)
}

// getDefaultHeader returns default headers
func getDefaultHeader() http.Header {
	return http.Header{"Content-Type": []string{"application/json"}}
}

// getAdditionalHeaders returns default header with addition correlation ID header
func getAdditionalHeaders() http.Header {
	return http.Header{"Content-Type": []string{"application/json"}, "Correlation-Id": []string{"some-id"}}
}

// getEmptyErrorManifest returns an empty manifest
func getEmptyErrorManifest() []reply.ErrorManifest {
	return []reply.ErrorManifest{}
}

// getDefaultErrorManifest returns the default manifest
func getDefaultErrorManifest() []reply.ErrorManifest {
	return []reply.ErrorManifest{
		{"example-404-error": reply.ErrorManifestItem{Title: "Resource Not Found", StatusCode: http.StatusNotFound}},
		{"example-name-validation-error": reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"}},
		{"example-dob-validation-error": reply.ErrorManifestItem{Title: "Validation Error", Detail: "Check your DoB, and try again.", Code: "100YT", StatusCode: http.StatusBadRequest}},
	}
}

// getMultiErrors returns example errors
func getMultiErrors() []error {
	return []error{
		errors.New("example-dob-validation-error"),
		errors.New("example-name-validation-error"),
	}
}

// getMultiErrorsWithMissingErr returns example errors with one error
// that does not exist in manifest
func getMultiErrorsWithMissingErr() []error {
	return append(getMultiErrors(), errors.New("example-missing-error"))
}

// getExampleErrorOne returns example error (1)
func getExampleErrorOne() error {
	return errors.New("example-404-error")
}

// getBlankResponseBody returns default blank response body
func getBlankResponseBody() string {
	return `{"data":"{}"}`
}

// getBlankResponseWithMetaBody returns default blank response body with meta-data
func getBlankResponseWithMetaBody() string {
	return `{"data":"{}","meta":{"example":"meta in response"}}`
}

// getDataResponseBody returns test data response body
func getDataResponseBody() string {
	return `{"data":{"id":"some-id","name":"john doe"}}`
}

// getDataResponseWithMetaBody returns test data response body with meta-data
func getDataResponseWithMetaBody() string {
	return `{"data":{"id":"some-id","name":"john doe"},"meta":{"example":"meta in response"}}`
}

// getFullTokenResponseBody returns test full token response body
func getFullTokenResponseBody() string {
	return `{"access_token":"test-token-1","refresh_token":"test-token-2"}`
}

// getSingleTokenResponseBody returns test single token response body
func getSingleTokenResponseBody() string {
	return `{"access_token":"test-token-1"}`
}

// getFullTokenResponseWithMetaBody returns test  full token response body with meta-data
func getFullTokenResponseWithMetaBody() string {
	return `{"access_token":"test-token-1","refresh_token":"test-token-2","meta":{"example":"meta in response"}}`
}

// getErrorResponseBody returns internal server error response body
func getErrorResponseISEBody() string {
	return `{"errors":[{"title":"Internal Server Error","status":"500"}]}`
}

// getErrorResponseForExampleErrorOne returns test error response body for getErrorResponseForExampleErrorOne function
func getErrorResponseForExampleErrorOne() string {
	return `{"errors":[{"title":"Resource Not Found","status":"404"}]}`
}

// getMultiErrorResponseMultiErrors returns test error response body for getMultiErrors function
func getMultiErrorResponseMultiErrors() string {
	return `{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}]}`
}

// getErrorResponseForExampleErrorOneWithMetaBody returns test error response body for getErrorResponseForExampleErrorOne function with meta-data
func getErrorResponseForExampleErrorOneWithMetaBody() string {
	return `{"errors":[{"title":"Resource Not Found","status":"404"}],"meta":{"example":"meta in response"}}`
}

// getMultiErrorResponseMultiErrorsWithMetaBody returns test error response body for getMultiErrors function with meta-data
func getMultiErrorResponseMultiErrorsWithMetaBody() string {
	return `{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}],"meta":{"example":"meta in response"}}`
}

// getTestUser returns user used by tests
func getTestUser() user {
	return user{
		ID:   "some-id",
		Name: "john doe",
	}
}

// getReplyFormattedHeader returns header in expected format for reply library
func getReplyFormattedHeader() map[string]string {
	return map[string]string{"correlation-id": "some-id"}
}

// getReplyFormattedMeta returns meta-data in expected format for reply library
func getReplyFormattedMeta() map[string]interface{} {
	return map[string]interface{}{
		"example": "meta in response",
	}
}
