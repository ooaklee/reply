package reply_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ooaklee/reply/v2"
	"github.com/stretchr/testify/require"
)

// misleadingStatusError checks that serialization cannot downgrade HTTP status.
type misleadingStatusError struct{ barError }

// GetStatusCode deliberately disagrees with the manifest for this regression.
func (*misleadingStatusError) GetStatusCode() string { return "201" }

// RefreshTransferObject returns independent custom state for each error.
func (*misleadingStatusError) RefreshTransferObject() reply.TransferObjectError {
	return &misleadingStatusError{}
}

// TestManifestControlsHTTPStatus retains manifest authority with custom errors.
func TestManifestControlsHTTPStatus(t *testing.T) {
	for _, status := range []int{0, 422, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			known := errors.New("known")
			r := reply.NewReplier([]reply.ErrorManifest{{known: {StatusCode: status}}}, reply.WithTransferObjectError(&misleadingStatusError{}))
			rec := httptest.NewRecorder()
			require.NoError(t, r.NewHTTPErrorResponse(rec, known))
			want := status
			if want == 0 {
				want = 400
			}
			require.Equal(t, want, rec.Code)
		})
	}
}

// invalidObject simulates custom response factories violating ownership.
type invalidObject struct {
	// TransferObject supplies unused setters; invalid factories must fail first.
	reply.TransferObject
	// mode chooses a nil, typed-nil or shared prototype result.
	mode string
}

// RefreshTransferObject returns the configured invalid object.
func (o *invalidObject) RefreshTransferObject() reply.TransferObject {
	switch o.mode {
	case "self":
		return o
	case "typed nil":
		return (*invalidObject)(nil)
	default:
		return nil
	}
}

// invalidErrorObject simulates invalid factories for individual public errors.
type invalidErrorObject struct {
	// TransferObjectError supplies unused methods which must not be reached.
	reply.TransferObjectError
	// mode selects the invalid ownership result.
	mode string
}

// RefreshTransferObject returns the configured invalid error object.
func (o *invalidErrorObject) RefreshTransferObject() reply.TransferObjectError {
	switch o.mode {
	case "self":
		return o
	case "typed nil":
		return (*invalidErrorObject)(nil)
	default:
		return nil
	}
}

// TestInvalidResponseInputs rejects malformed configuration without panic or a
// committed response. A nil Error in an otherwise valid request stays compatible.
func TestInvalidResponseInputs(t *testing.T) {
	known := errors.New("known")
	for _, tc := range []struct {
		name       string
		nilReplier bool
		nilRequest bool
		writer     http.ResponseWriter
		options    []reply.Option
		valid      bool
	}{
		{name: "nil replier", nilReplier: true, writer: httptest.NewRecorder()},
		{name: "nil request", nilRequest: true, writer: httptest.NewRecorder()},
		{name: "nil writer"},
		{name: "typed nil writer", writer: (*httptest.ResponseRecorder)(nil)},
		{name: "nil prototype", writer: httptest.NewRecorder(), options: []reply.Option{reply.WithTransferObject(nil)}},
		{name: "typed nil prototype", writer: httptest.NewRecorder(), options: []reply.Option{reply.WithTransferObject((*invalidObject)(nil))}},
		{name: "nil error prototype", writer: httptest.NewRecorder(), options: []reply.Option{reply.WithTransferObjectError(nil)}},
		{name: "nil option", writer: httptest.NewRecorder(), options: []reply.Option{nil}, valid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := reply.NewReplier([]reply.ErrorManifest{{known: {StatusCode: 409}}}, tc.options...)
			if tc.nilReplier {
				r = nil
			}
			req := &reply.NewResponseRequest{Writer: tc.writer, Error: known}
			if tc.nilRequest {
				req = nil
			}
			err := r.NewHTTPResponse(req)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if rec, ok := tc.writer.(*httptest.ResponseRecorder); ok && rec != nil {
					require.Empty(t, rec.Body.String())
				}
			}
		})
	}
	for _, kind := range []string{"response", "error"} {
		t.Run(kind, func(t *testing.T) {
			for _, mode := range []string{"nil", "typed nil", "self"} {
				t.Run(mode, func(t *testing.T) {
					var option reply.Option
					if kind == "response" {
						option = reply.WithTransferObject(&invalidObject{mode: mode})
					} else {
						option = reply.WithTransferObjectError(&invalidErrorObject{mode: mode})
					}
					r := reply.NewReplier([]reply.ErrorManifest{{known: {StatusCode: 409}}}, option)
					rec := httptest.NewRecorder()
					require.Error(t, r.NewHTTPErrorResponse(rec, known))
					require.Empty(t, rec.Body.String())
				})
			}
		})
	}
}

// TestOptionalContextAndAttributes retains nil-attribute and blank contracts;
// an unknown response emits one safe observation with a non-nil context.
func TestOptionalContextAndAttributes(t *testing.T) {
	for _, tc := range []struct {
		name         string
		err          error
		want         int
		observations int
	}{
		{"nil error remains blank", nil, 200, 0},
		{"unknown error", errors.New("private"), 500, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			observed := 0
			r := reply.NewReplier(nil, reply.WithUnmappedErrorObserver(func(ctx context.Context) { require.NotNil(t, ctx); observed++ }))
			rec := httptest.NewRecorder()
			require.NoError(t, r.NewHTTPErrorResponse(rec, tc.err, nil, reply.WithContext(nil)))
			require.Equal(t, tc.want, rec.Code)
			require.Equal(t, tc.observations, observed)
		})
	}
}
