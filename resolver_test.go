package reply_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ooaklee/reply/v2"
	"github.com/stretchr/testify/require"
)

// classifiedError separates a public classification from an internal cause.
type classifiedError struct {
	// keys are the deliberately explicit Is matches for this fixture node.
	keys []error
	// cause remains available to retry/cause inspection, never response text.
	cause error
}

// Error panics so any accidental text-based resolution fails the test.
func (e *classifiedError) Error() string { panic("diagnostic text must not be inspected") }

// Is compares this node only; it does not search joined branches.
func (e *classifiedError) Is(target error) bool {
	for _, key := range e.keys {
		if target == key {
			return true
		}
	}
	return false
}

// Unwrap exposes the fixture's private cause.
func (e *classifiedError) Unwrap() error { return e.cause }

// malformedJoin permits empty and nil-containing joins that errors.Join removes.
type malformedJoin []error

// Error panics to prevent reliance on text or newline heuristics.
func (malformedJoin) Error() string { panic("diagnostic text must not be inspected") }

// Unwrap returns each independent branch unchanged.
func (e malformedJoin) Unwrap() []error { return e }

// cyclicError verifies the traversal budget also bounds cycles.
type cyclicError struct{}

// Error is intentionally not suitable for classification.
func (*cyclicError) Error() string { return "cycle" }

// Unwrap returns this node repeatedly.
func (e *cyclicError) Unwrap() error { return e }

// dynamicError has a comparable type but may hold an uncomparable value.
type dynamicError struct{ value interface{} }

// Error is diagnostic text, not an identity.
func (dynamicError) Error() string { return "dynamic" }

// wrappedNodes constructs a precise total-node fixture.
func wrappedNodes(err error, count int) error {
	for i := 1; i < count; i++ {
		err = fmt.Errorf("private diagnostic: %w", err)
	}
	return err
}

// repeatedErrors constructs a breadth-boundary fixture.
func repeatedErrors(err error, count int) []error {
	result := make([]error, count)
	for i := range result {
		result[i] = err
	}
	return result
}

// TestResolveErrors pins complete structural classification across all roots.
func TestResolveErrors(t *testing.T) {
	first, second, private := errors.New("first"), errors.New("second"), errors.New("private")
	classified := &classifiedError{keys: []error{first}, cause: private}
	exact := &classifiedError{keys: []error{first, second}, cause: private}
	var typedNil *classifiedError
	r := reply.NewReplier([]reply.ErrorManifest{{first: {Code: "FIRST"}, second: {Code: "SECOND"}, exact: {Code: "EXACT"}, nil: {Code: "INVALID"}, typedNil: {Code: "INVALID"}}})
	for _, tc := range []struct {
		name string
		errs []error
		want []error
	}{
		{"direct", []error{first}, []error{first}},
		{"wrapped", []error{fmt.Errorf("private: %w", first)}, []error{first}},
		{"multiline wrapper", []error{fmt.Errorf("private\n%w", first)}, []error{first}},
		{"singleton join", []error{errors.Join(first)}, []error{first}},
		{"ordered nested join", []error{errors.Join(second, errors.Join(first, second))}, []error{second, first}},
		{"explicit roots", []error{second, first, second}, []error{second, first}},
		{"custom classification", []error{classified}, []error{first}},
		{"exact before custom Is", []error{exact}, []error{exact}},
		{"classification and unknown sibling", []error{malformedJoin{classified, private}}, nil},
		{"ambiguous Is", []error{&classifiedError{keys: []error{first, second}}}, nil},
		{"unknown", []error{private}, nil},
		{"unknown sibling first", []error{private, first}, nil},
		{"unknown sibling last", []error{first, private}, nil},
		{"text lookalike", []error{errors.New("first")}, nil},
		{"newlines are not structure", []error{errors.New("first\nsecond")}, nil},
		{"no roots", nil, nil},
		{"nil root", []error{nil}, nil},
		{"typed nil root", []error{typedNil}, nil},
		{"empty join", []error{malformedJoin{}}, nil},
		{"nil child", []error{malformedJoin{first, nil}}, nil},
		{"dynamic incomparability", []error{dynamicError{value: []int{1}}}, nil},
		{"cycle", []error{&cyclicError{}}, nil},
		{"63 nodes", []error{wrappedNodes(first, 63)}, []error{first}},
		{"64 nodes", []error{wrappedNodes(first, 64)}, []error{first}},
		{"65 nodes", []error{wrappedNodes(first, 65)}, nil},
		{"wide within budget", []error{malformedJoin(repeatedErrors(first, 63))}, []error{first}},
		{"wide over budget", []error{malformedJoin(repeatedErrors(first, 64))}, nil},
		{"64 roots", repeatedErrors(first, 64), []error{first}},
		{"65 roots", repeatedErrors(first, 65), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := r.ResolveErrors(tc.errs...)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.want != nil, ok)
		})
	}
}

// TestResolvedResponseContract checks public statuses, ordering and override
// behavior on both response entry points. Cases are serial for log capture.
func TestResolvedResponseContract(t *testing.T) {
	first, second, server := errors.New("first"), errors.New("second"), errors.New("server")
	for _, tc := range []struct {
		name      string
		errs      []error
		status    int
		codes     []string
		unmatched bool
	}{
		{"direct override", []error{first}, 422, []string{"OVERRIDE"}, false},
		{"wrapped override", []error{fmt.Errorf("private: %w", first)}, 422, []string{"OVERRIDE"}, false},
		{"joined order", []error{errors.Join(second, first)}, 409, []string{"SECOND", "OVERRIDE"}, false},
		{"server dominance", []error{errors.Join(first, server)}, 503, []string{"SERVER"}, false},
		{"server cannot hide unknown", []error{server, errors.New("private")}, 500, []string{""}, true},
		{"unknown", []error{errors.New("private")}, 500, []string{""}, true},
		{"typed nil", []error{(*classifiedError)(nil)}, 500, []string{""}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"Error", "Errors"} {
				t.Run(mode, func(t *testing.T) {
					var logs bytes.Buffer
					previous := log.Writer()
					log.SetOutput(&logs)
					t.Cleanup(func() { log.SetOutput(previous) })
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					observations := 0
					r := reply.NewReplier([]reply.ErrorManifest{{first: {Code: "FIRST", StatusCode: 400}, second: {Code: "SECOND", StatusCode: 409}, server: {Code: "SERVER", StatusCode: 503}}, {first: {Code: "OVERRIDE", StatusCode: 422}}}, reply.WithUnmappedErrorObserver(func(got context.Context) { require.Same(t, ctx, got); observations++ }))
					rec := httptest.NewRecorder()
					var err error
					if mode == "Error" {
						err = r.NewHTTPErrorResponse(rec, errors.Join(tc.errs...), reply.WithContext(ctx))
					} else {
						err = r.NewHTTPMultiErrorResponse(rec, tc.errs, reply.WithContext(ctx))
					}
					require.NoError(t, err)
					require.Equal(t, tc.status, rec.Code)
					var body struct {
						Errors []struct {
							Code string `json:"code"`
						} `json:"errors"`
					}
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
					var codes []string
					for _, item := range body.Errors {
						codes = append(codes, item.Code)
					}
					require.Equal(t, tc.codes, codes)
					require.Equal(t, tc.unmatched, observations == 1)
					require.NotContains(t, rec.Body.String(), "private")
					require.Empty(t, logs.String())
				})
			}
		})
	}
}

// TestSharedReplier isolates default and custom response state during parallel
// mixed requests. Prototypes are deliberately shared; each writer is case-owned.
func TestSharedReplier(t *testing.T) {
	for _, custom := range []bool{false, true} {
		custom := custom
		t.Run(fmt.Sprint("custom=", custom), func(t *testing.T) {
			known := errors.New("known")
			opts := []reply.Option{}
			if custom {
				opts = append(opts, reply.WithTransferObject(&fooReplyTransferObject{}), reply.WithTransferObjectError(&barError{}))
			}
			r := reply.NewReplier([]reply.ErrorManifest{{known: {Code: "KNOWN", StatusCode: 409}}}, opts...)
			for i := 0; i < 64; i++ {
				i := i
				t.Run(fmt.Sprint(i), func(t *testing.T) {
					t.Parallel()
					want, got := httptest.NewRecorder(), httptest.NewRecorder()
					request := &reply.NewResponseRequest{Writer: want, Headers: map[string]string{"X-Request": fmt.Sprint(i)}, Meta: map[string]interface{}{"request": i}}
					switch i % 4 {
					case 0:
						request.Error = fmt.Errorf("private: %w", known)
					case 1:
						request.TokenOne = fmt.Sprint(i)
					case 2:
						request.Data = map[string]int{"request": i}
					case 3:
						request.StatusCode = 202
					}
					require.NoError(t, reply.NewReplier([]reply.ErrorManifest{{known: {Code: "KNOWN", StatusCode: 409}}}, opts...).NewHTTPResponse(request))
					request.Writer = got
					require.NoError(t, r.NewHTTPResponse(request))
					require.Equal(t, want.Code, got.Code)
					require.Equal(t, want.Header(), got.Header())
					require.Equal(t, want.Body.String(), got.Body.String())
				})
			}
		})
	}
}

// TestManifestSnapshot verifies entries and map ownership freeze at construction;
// nested metadata is intentionally caller-owned immutable data, not deep-copied.
func TestManifestSnapshot(t *testing.T) {
	for _, status := range []int{0, 202, 400, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			known := errors.New("known")
			source := reply.ErrorManifest{known: {Title: "Public", Code: "CODE", StatusCode: status}}
			r := reply.NewReplier([]reply.ErrorManifest{source})
			delete(source, known)
			source[known] = reply.ErrorManifestItem{Code: "MUTATED", StatusCode: 418}
			rec := httptest.NewRecorder()
			require.NoError(t, r.NewHTTPErrorResponse(rec, known))
			require.Contains(t, rec.Body.String(), "CODE")
			require.NotContains(t, rec.Body.String(), "MUTATED")
			want := status
			if want == 0 {
				want = 400
			}
			require.Equal(t, want, rec.Code)
		})
	}
}

// failingWriter preserves a native transport failure for errors.Is inspection.
type failingWriter struct {
	http.ResponseWriter
	cause error
}

// Write reports a disconnected response sink without storing the body.
func (w failingWriter) Write([]byte) (int, error) { return 0, w.cause }

// TestWriteFailureCauses verifies transport and JSON causes survive wrapping.
func TestWriteFailureCauses(t *testing.T) {
	for _, kind := range []string{"writer", "encoding"} {
		t.Run(kind, func(t *testing.T) {
			r := reply.NewReplier(nil)
			rec := httptest.NewRecorder()
			cause := errors.New("disconnected")
			if kind == "writer" {
				err := r.NewHTTPDataResponse(failingWriter{rec, cause}, 200, "data")
				require.ErrorIs(t, err, cause)
			} else {
				err := r.NewHTTPDataResponse(rec, 200, make(chan int))
				var unsupported *json.UnsupportedTypeError
				require.ErrorAs(t, err, &unsupported)
			}
		})
	}
}

// BenchmarkResolveErrors compares canonical and structural request-time lookup
// using one premerged immutable manifest; it does not benchmark HTTP encoding.
func BenchmarkResolveErrors(b *testing.B) {
	known := errors.New("known")
	r := reply.NewReplier([]reply.ErrorManifest{{known: {Code: "KNOWN"}}})
	for _, tc := range []struct {
		name string
		err  error
	}{{"exact", known}, {"wrapped", fmt.Errorf("context: %w", known)}, {"joined", errors.Join(known, known)}} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, ok := r.ResolveErrors(tc.err); !ok {
					b.Fatal("unresolved")
				}
			}
		})
	}
}
