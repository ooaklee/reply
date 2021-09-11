// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ooaklee/reply"
	"github.com/stretchr/testify/assert"
)

type baseTestResponse struct {
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type baseStatusMessageResponse struct {
	Status struct {
		Message string `json:"message,omitempty"`
	} `json:"status,omitempty"`
}

type user struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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

				assert.Equal(t, stringWithNewLine(`{"data":"{}"}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"data":"{}"}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Headers:    map[string]string{"correlation-id": "some-id"},
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":"{}"}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":"{}","meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Blank response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 302,
				Headers:    map[string]string{"correlation-id": "some-id"},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: 302,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":"{}","meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		//////// Data Response ////////
		{
			name:      "Success - Data response",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":{"id":"some-id","name":"john doe"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				StatusCode: 201,
				Headers:    map[string]string{"correlation-id": "some-id"},
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":{"id":"some-id","name":"john doe"}}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 201,
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":{"id":"some-id","name":"john doe"},"meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 201,
				Headers:    map[string]string{"correlation-id": "some-id"},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
			},
			expectedStatusCode: 201,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":{"id":"some-id","name":"john doe"},"meta":{"example":"meta in response"}}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"access_token":"test-token-1","refresh_token":"test-token-2"}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"access_token":"test-token-1"}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Additional Headers",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
				StatusCode: 200,
				Headers:    map[string]string{"correlation-id": "some-id"},
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"access_token":"test-token-1","refresh_token":"test-token-2"}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 200,
				TokenOne:   "test-token-1",
				TokenTwo:   "test-token-2",
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"access_token":"test-token-1","refresh_token":"test-token-2","meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response with Meta-Information & Additional Header",
			manifests: getEmptyErrorManifest(),
			request: reply.NewResponseRequest{
				StatusCode: 200,
				Headers:    map[string]string{"correlation-id": "some-id"},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
				TokenOne: "test-token-1",
				TokenTwo: "test-token-2",
			},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"access_token":"test-token-1","refresh_token":"test-token-2","meta":{"example":"meta in response"}}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Internal Server Error","status":"500"}]}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Resource Not Found","status":"404"}]}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}]}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Internal Server Error","status":"500"}]}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Additional Headers",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error:   getExampleErrorOne(),
				Headers: map[string]string{"correlation-id": "some-id"},
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Resource Not Found","status":"404"}]}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Additional Headers",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors:  getMultiErrors(),
				Headers: map[string]string{"correlation-id": "some-id"},
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}]}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Meta-Information",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error: getExampleErrorOne(),
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: 404,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Resource Not Found","status":"404"}],"meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Meta-Information",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors: getMultiErrors(),
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}],"meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response with Meta-Information & Additional Header",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error:   getExampleErrorOne(),
				Headers: map[string]string{"correlation-id": "some-id"},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Resource Not Found","status":"404"}],"meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		{
			name:      "Success - Multi error response with Meta-Information & Additional Header",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors:  getMultiErrors(),
				Headers: map[string]string{"correlation-id": "some-id"},
				Meta: map[string]interface{}{
					"example": "meta in response",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}],"meta":{"example":"meta in response"}}`), returnedBody)

				assert.Equal(t, getAdditionalHeaders(), w.Header())
			},
		},
		///////////////////////////////
		////// Response Ranking ///////
		{
			name:      "Success - Multi Error response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Errors: getMultiErrors(),
				Error:  getExampleErrorOne(),
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusBadRequest,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Validation Error","detail":"Check your DoB, and try again.","status":"400","code":"100YT"},{"title":"Validation Error","detail":"The name provided does not meet validation requirements","about":"www.example.com/reply/validation/1011","status":"400","code":"1011"}]}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Error response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Error: getExampleErrorOne(),
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"errors":[{"title":"Resource Not Found","status":"404"}]}`), returnedBody)

				assert.Equal(t, getDefaultHeader(), w.Header())
			},
		},
		{
			name:      "Success - Data response should take precedence",
			manifests: getDefaultErrorManifest(),
			request: reply.NewResponseRequest{
				Data: user{
					ID:   "some-id",
					Name: "john doe",
				},
				StatusCode: 201,
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				returnedBody := w.Body.String()

				assert.Equal(t, stringWithNewLine(`{"data":{"id":"some-id","name":"john doe"}}`), returnedBody)

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

				assert.Equal(t, stringWithNewLine(`{"data":"{}"}`), returnedBody)

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

func TestReplier_AideNewHTTPErrorResponse(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		err                error
		StatusCode         int
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name: "Success - Resource not found",
			manifests: append([]reply.ErrorManifest{
				{"test-404-error": reply.ErrorManifestItem{Title: "resource not found", StatusCode: http.StatusNotFound}},
			},
				reply.ErrorManifest{
					"test-401-error": reply.ErrorManifestItem{Title: "unauthorized", StatusCode: http.StatusUnauthorized},
				},
			),
			err:                errors.New("test-404-error"),
			expectedStatusCode: http.StatusNotFound,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseStatusMessageResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				expectedResponse := baseStatusMessageResponse{}
				expectedResponse.Status.Message = "resource not found"
				assert.Equal(t, expectedResponse, response)
			},
		},
		{
			name:               "Failure - Error not in manifest",
			manifests:          []reply.ErrorManifest{},
			err:                errors.New("test-404-error"),
			expectedStatusCode: http.StatusInternalServerError,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseStatusMessageResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				expectedResponse := baseStatusMessageResponse{}
				expectedResponse.Status.Message = "Internal Server Error"
				assert.Equal(t, expectedResponse, response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			replier := reply.NewReplier(test.manifests)

			replier.NewHTTPErrorResponse(w, test.err)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_AideNewHTTPDataResponse(t *testing.T) {

	type user struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		data               interface{}
		StatusCode         int
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:       "Success - Created Mock user",
			manifests:  []reply.ErrorManifest{},
			StatusCode: 201,
			data: user{
				ID:   "new-uuid",
				Name: "Test User",
			},
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseTestResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				assert.Equal(t, baseTestResponse{Data: map[string]interface{}{"id": "new-uuid", "name": "Test User"}}, response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			replier := reply.NewReplier(test.manifests)

			replier.NewHTTPDataResponse(w, test.StatusCode, test.data)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_AideNewHTTPTokenResponse(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		accessToken        string
		refreshToken       string
		StatusCode         int
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Success - Access Token response",
			manifests:          []reply.ErrorManifest{},
			accessToken:        "test-access-token",
			StatusCode:         200,
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseTestResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				assert.Equal(t, baseTestResponse{AccessToken: "test-access-token"}, response)
			},
		},
		{
			name:               "Success - Full Token response",
			manifests:          []reply.ErrorManifest{},
			accessToken:        "test-access-token",
			refreshToken:       "test-refresh-token",
			StatusCode:         200,
			expectedStatusCode: http.StatusOK,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseTestResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				assert.Equal(t, baseTestResponse{AccessToken: "test-access-token", RefreshToken: "test-refresh-token"}, response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			replier := reply.NewReplier(test.manifests)

			replier.NewHTTPTokenResponse(w, test.StatusCode, test.accessToken, test.refreshToken)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_AideNewHTTPBlankResponse(t *testing.T) {

	tests := []struct {
		name               string
		manifests          []reply.ErrorManifest
		StatusCode         int
		assertResponse     func(w *httptest.ResponseRecorder, t *testing.T)
		expectedStatusCode int
	}{
		{
			name:               "Success",
			manifests:          []reply.ErrorManifest{},
			StatusCode:         201,
			expectedStatusCode: http.StatusCreated,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseTestResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				assert.Equal(t, baseTestResponse{Data: "{}"}, response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			replier := reply.NewReplier(test.manifests)

			replier.NewHTTPBlankResponse(w, test.StatusCode)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

// unmarshalResponseBody handles unmarshalling recorder's response to specified
// response body
func unmarshalResponseBody(w *httptest.ResponseRecorder, responseBody interface{}) error {
	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, responseBody); err != nil {
		return err
	}

	return nil
}
