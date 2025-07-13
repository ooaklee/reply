// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			name:               "Success - Wrapped Error Response (Multi Error response)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getWrappedErrors(),
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrors()), returnedBody)

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
		///////////////////////////////
		/////// With Custom TOE ///////
		{
			name:               "Failure - Error response (Using custom TOE)",
			manifests:          getEmptyErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Wrapped Error Response (Multi Error response w/ Custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getWrappedErrors(),
			transferObjecError: &barError{},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Additional Headers (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneUsingCustomTOE()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Meta-Information (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 404,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Meta-Information & Additional Header (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		///// With Custom TOE & TO ////
		{
			name:               "Failure - Error response (Using custom TOE & TO)",
			manifests:          getEmptyErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBodyUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Wrapped Error Response (Multi Error response w/ Custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getWrappedErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Additional Headers (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Meta-Information (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 404,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Error response with Meta-Information & Additional Header (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedError:        getExampleErrorOne(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOEAndTO()), returnedBody)

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
		///////////////////////////////
		/////// With Custom TOE ///////
		{
			name:               "Failure - Multi Error response (Using custom TOE)",
			manifests:          getEmptyErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Additional Headers (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOE()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Meta-Information (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 400,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Meta-Information & Additional Header (Using custom TOE)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOE()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		///// With Custom TOE & TO ////
		{
			name:               "Failure - Multi Error response (Using custom TOE & TO)",
			manifests:          getEmptyErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getErrorResponseISEBodyUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Additional Headers (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Meta-Information (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
			},
			expectedStatusCode: 400,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOEAndTO()), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:               "Success - Multi Error response with Meta-Information & Additional Header (Using custom TOE & TO)",
			manifests:          getDefaultErrorManifest(),
			passedErrors:       getMultiErrors(),
			transferObjecError: &barError{},
			transferObject:     &fooReplyTransferObject{},
			responseAttributes: []reply.ResponseAttributes{
				reply.WithMeta(getReplyFormattedMeta()),
				reply.WithHeaders(getReplyFormattedHeader()),
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOEAndTO()), returnedBody)

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
		{"example-dob-validation-error": reply.ErrorManifestItem{Title: "Validation Error", Detail: "Check your DoB, and try again.", Code: "100YT", StatusCode: http.StatusBadRequest}},
		{"example-name-validation-error": reply.ErrorManifestItem{Title: "Validation Error", Detail: "The name provided does not meet validation requirements", StatusCode: http.StatusBadRequest, About: "www.example.com/reply/validation/1011", Code: "1011"}},
	}
}

// getMultiErrors returns example errors
func getMultiErrors() []error {
	return []error{
		errors.New("example-dob-validation-error"),
		errors.New("example-name-validation-error"),
	}
}

// getWrappedErrors returns wrapped errors
func getWrappedErrors() error {

	err := errors.New("example-name-validation-error")
	errTwo := errors.New("example-dob-validation-error")

	return errors.Join(err, errTwo)

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

////////////////////////
//// For custom TOE ////

// getErrorResponseBodyUsingCustomTOE returns internal server error response body
// based on custom Transfer Object Error
func getErrorResponseISEBodyUsingCustomTOE() string {
	return `{"errors":[{"title":"Internal Server Error","more":{"status":"500"}}]}`
}

// getErrorResponseForExampleErrorOneUsingCustomTOE returns test error response body for getErrorResponseForExampleErrorOne function
// based on custom Transfer Object Error
func getErrorResponseForExampleErrorOneUsingCustomTOE() string {
	return `{"errors":[{"title":"Resource Not Found","more":{"status":"404"}}]}`
}

// getMultiErrorResponseMultiErrorsUsingCustomTOE returns test error response body for getMultiErrors function
// based on custom Transfer Object Error
func getMultiErrorResponseMultiErrorsUsingCustomTOE() string {
	return `{"errors":[{"title":"Validation Error","message":"Check your DoB, and try again.","more":{"status":"400","code":"100YT"}},{"title":"Validation Error","message":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","more":{"status":"400","code":"1011"}}]}`
}

// getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOE returns test error response body for getErrorResponseForExampleErrorOne function with meta-data
// based on custom Transfer Object Error
func getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOE() string {
	return `{"errors":[{"title":"Resource Not Found","more":{"status":"404"}}],"meta":{"example":"meta in response"}}`
}

// getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOE returns test error response body for getMultiErrors function with meta-data
// based on custom Transfer Object Error
func getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOE() string {
	return `{"errors":[{"title":"Validation Error","message":"Check your DoB, and try again.","more":{"status":"400","code":"100YT"}},{"title":"Validation Error","message":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","more":{"status":"400","code":"1011"}}],"meta":{"example":"meta in response"}}`
}

////////////////////////
// For custom TOE & TO /

// getErrorResponseBodyUsingCustomTOEAndTO returns internal server error response body
// based on custom Transfer Object Error & Transfer Object
func getErrorResponseISEBodyUsingCustomTOEAndTO() string {
	return `{"bar":{"errors":[{"title":"Internal Server Error","more":{"status":"500"}}]}}`
}

// getErrorResponseForExampleErrorOneUsingCustomTOEAndTO returns test error response body for getErrorResponseForExampleErrorOne function
// based on custom Transfer Object Error & Transfer Object
func getErrorResponseForExampleErrorOneUsingCustomTOEAndTO() string {
	return `{"bar":{"errors":[{"title":"Resource Not Found","more":{"status":"404"}}]}}`
}

// getMultiErrorResponseMultiErrorsUsingCustomTOEAndTO returns test error response body for getMultiErrors function
// based on custom Transfer Object Error & Transfer Object
func getMultiErrorResponseMultiErrorsUsingCustomTOEAndTO() string {
	return `{"bar":{"errors":[{"title":"Validation Error","message":"Check your DoB, and try again.","more":{"status":"400","code":"100YT"}},{"title":"Validation Error","message":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","more":{"status":"400","code":"1011"}}]}}`
}

// getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOEAndTO returns test error response body for getErrorResponseForExampleErrorOne function with meta-data
// based on custom Transfer Object Error & Transfer Object
func getErrorResponseForExampleErrorOneWithMetaBodyUsingCustomTOEAndTO() string {
	return `{"bar":{"errors":[{"title":"Resource Not Found","more":{"status":"404"}}],"meta":{"example":"meta in response"}}}`
}

// getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOEAndTO returns test error response body for getMultiErrors function with meta-data
// based on custom Transfer Object Error & Transfer Object
func getMultiErrorResponseMultiErrorsWithMetaBodyUsingCustomTOEAndTO() string {
	return `{"bar":{"errors":[{"title":"Validation Error","message":"Check your DoB, and try again.","more":{"status":"400","code":"100YT"}},{"title":"Validation Error","message":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","more":{"status":"400","code":"1011"}}],"meta":{"example":"meta in response"}}}`
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

// ///////////////////////////////////////////////
// ///// Custom Transition Object Example ////////
// This is an example of how you can create a
// custom response structure based on your
// requirements.
//
// Sourced from `examples/example_simple_api.go`
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
//
// Sourced from `examples/example_simple_api.go`

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
