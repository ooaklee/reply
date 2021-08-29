// Copyright (C) 2021 by Leon Silcott <leon@boasi.io>. All rights reserved.
// Use of this source code is governed under MIT License.
// See the [LICENSE](https://github.com/ooaklee/reply/blob/master/LICENSE) for details.

package reply_test

import (
	"encoding/json"
	"errors"
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

func TestReplier_NewHTTPResponseForTokens(t *testing.T) {

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
		{
			name:               "Failed - Default response sent",
			manifests:          []reply.ErrorManifest{},
			expectedStatusCode: 200,
			assertResponse: func(w *httptest.ResponseRecorder, t *testing.T) {

				response := baseTestResponse{}

				err := unmarshalResponseBody(w, &response)
				if err != nil {
					t.Fatalf("cannot get response content: %v", err)
				}

				expectedResponse := baseTestResponse{Data: "{}"}

				assert.Equal(t, expectedResponse, response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			replier := reply.NewReplier(test.manifests)

			replier.NewHTTPResponse(&reply.NewResponseRequest{
				Writer:       w,
				StatusCode:   test.StatusCode,
				AccessToken:  test.accessToken,
				RefreshToken: test.refreshToken,
			})

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPResponseForData(t *testing.T) {

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

			replier.NewHTTPResponse(&reply.NewResponseRequest{
				Writer:     w,
				StatusCode: test.StatusCode,
				Data:       test.data,
			})

			assert.Equal(t, test.expectedStatusCode, w.Code)
			test.assertResponse(w, t)
		})
	}
}

func TestReplier_NewHTTPResponseForError(t *testing.T) {

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
			manifests: []reply.ErrorManifest{
				{"test-404-error": reply.ErrorManifestItem{Message: "resource not found", StatusCode: http.StatusNotFound}},
			},
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

			replier.NewHTTPResponse(&reply.NewResponseRequest{
				Writer: w,
				Error:  test.err,
			})

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
