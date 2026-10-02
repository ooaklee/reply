# Integrating the response-resolution upgrade

This unreleased change is developed on `feat/improve-functionality`. Keep the
`github.com/ooaklee/reply/v2` module path. The response helper signatures and
custom transfer-object interfaces remain available.

## Pinning during integration

Resolve the branch with Go tooling; `go.mod` records an immutable pseudo-version
for its current commit, not a moving branch reference:

```sh
go get github.com/ooaklee/reply/v2@feat/improve-functionality
go test -race ./...
```

For reproducible adoption across services, select the same verified commit in
each module with `go get github.com/ooaklee/reply/v2@<commit>`. Commit both module
files. Update that pin deliberately after subsequent branch changes; do not
invent pseudo-versions or disable module checksum verification.

For unpublished local work, a temporary sibling replacement may be used:

```go
replace github.com/ooaklee/reply/v2 => ../reply
```

Replacements in dependencies are ignored. Each top-level host or workspace must
configure its own replacement. Remove it and rerun integration tests against the
published commit before relying on standalone builds or CI.

## Error resolution

All error-response entry points now resolve errors structurally against the
merged manifest. Expected failures still require entries; the resolver does not
infer public details from error text.

- Direct sentinels retain their existing wire fields and status.
- `%w` wrappers, single-element joins and nested joins resolve correctly.
- Joins and explicit error slices use first-seen input order and deduplicate
  repeated identities. Previously, joins followed unpredictable map iteration;
  explicit slices could repeat the same error. Clients should key on error codes.
- Every independent branch must resolve. Unknown, nil/typed-nil, malformed,
  ambiguous or over-budget input produces a generic 500 for the whole response.
  An unknown sibling is never discarded behind a recognized client/server error.
- An exact registered node classifies its own causes. Otherwise, a custom `Is`
  method may match exactly one manifest identity for that node. Custom methods
  must not recursively classify independent siblings and must terminate; reply
  does not recover panics or interrupt application-provided error methods.
  Unregistered nodes with `Unwrap() []error` always have their children checked;
  an `Is` method on the same node cannot hide an independent failure.
- Traversal is bounded to 64 nodes across all roots, not just 64 levels deep.
- Later manifest entries for the same identity still win. A mapped 5xx still
  dominates an otherwise fully mapped collection; otherwise its first error's
  status wins, defaulting to 400. Identity is not error-message equality.

`Replier.ResolveErrors` exposes the same resolution without HTTP effects. It
returns canonical keys and `true`, or `nil, false` for incomplete resolution.
For an error-only adapter, put a possibly nil error into an explicit `Errors`
slice to obtain the generic failure. A generic request with no `Error` and no
`Errors` intentionally retains the existing successful blank-response contract.

## Diagnostics and request context

**Behavior change:** unresolved errors are no longer printed to the standard
logger. Reply neither formats the original error nor includes its diagnostics in
the fallback response. To retain observability, opt in to a safe callback:

```go
replier := reply.NewReplier(manifests,
    reply.WithUnmappedErrorObserver(func(ctx context.Context) {
        // Increment a metric or record a redacted request identifier from ctx.
        // The callback receives no original error or private diagnostic string.
    }),
)

err := replier.NewHTTPErrorResponse(w, domainErr, reply.WithContext(r.Context()))
```

The observer runs once per unresolved response, synchronously. It must be
concurrency-safe, non-blocking and non-panicking. If context was omitted it
receives `context.Background()`. Context is never serialized and does not grant
authority, change status, or implement writer cancellation. Handle returned
write/encoding errors yourself; their native causes now survive via `%w`.
Manifest fields and caller metadata/headers remain your responsibility to redact.
The observer records resolution failure, not successful HTTP delivery; a later
factory, encoder or writer failure may still prevent the response from being sent.

## Shared repliers and custom transfer objects

Response state is now request-local. A replier using the default transfer objects
can be reused concurrently after construction. Its manifest entries are copied
and merged once; nested metadata values are **not deep-copied** and must remain
immutable. Do not apply options to a live replier or mutate shared response data.

Custom `RefreshTransferObject` methods must return independent state without
mutating their prototype and must be safe to call concurrently. Nil/typed-nil
prototypes, nil results and factories returning the same prototype pointer now
return errors instead of proceeding. Hidden aliasing cannot be detected: custom
implementations remain responsible for independent nested state. Existing
custom envelopes are retained; no new factory interface is required.

Nil options and response attributes are ignored. A nil response request, nil
writer or invalid prototype returns an error. Default token, data, metadata,
header and blank-response precedence and serialization remain unchanged.
HTTP error status comes from the manifest, not a custom error serializer's
`GetStatusCode` result. Custom serializers should faithfully represent that status.

## Verification

```sh
go test -race -count=1 ./...
go vet ./...
go test -run '^$' -bench BenchmarkResolveErrors -benchmem
```

The table-driven suite covers structural failures, stable status/order, private
diagnostics, prototype validation, custom envelopes, shared concurrent repliers,
manifest snapshots and native writer/encoder errors. Benchmarks measure lookup
with one premerged manifest; they are not end-to-end service latency claims.
