package reply

import "reflect"

// ResolveErrors resolves complete error trees to manifest keys in first-seen
// order, removing duplicate keys. It understands single-error wrappers, Go
// joins and custom Is methods without parsing or logging diagnostic strings.
// It returns false if any independent branch is unknown, malformed, ambiguous
// or exceeds the 64-node traversal budget; partial results are never returned.
// An empty input and nil errors are invalid, not successful responses.
//
// Exact registered nodes classify their own causes. Otherwise, a custom Is
// method may classify its node as exactly one manifest key; it must not search
// other joined branches. An unregistered multi-unwrapper is always traversed;
// its Is method cannot hide independent children. Is and Unwrap methods must terminate and be safe to
// call. Those user implementations are not sandboxed or recovered from panic.
// The replier's immutable, premerged manifest is reused across calls.
func (r *Replier) ResolveErrors(errs ...error) ([]error, bool) {
	if r == nil || len(errs) == 0 || len(errs) > 64 {
		return nil, false
	}
	// Most HTTP failures are one canonical sentinel. Avoid traversal bookkeeping
	// for that case while retaining typed-nil and comparability checks.
	if len(errs) == 1 && !isNilValue(errs[0]) && reflect.ValueOf(errs[0]).Comparable() {
		if _, exists := r.errorManifest[errs[0]]; exists {
			return []error{errs[0]}, true
		}
	}
	resolver := errorResolver{manifest: r.errorManifest, seen: make(map[error]struct{}), remaining: 64}
	for _, err := range errs {
		if !resolver.visit(err) {
			return nil, false
		}
	}
	return resolver.resolved, true
}

// errorResolver owns one bounded traversal and never mutates the manifest.
type errorResolver struct {
	// manifest contains stable registered identities and public response fields.
	manifest ErrorManifest
	// seen tracks emitted keys without skipping validation of independent branches.
	seen map[error]struct{}
	// resolved preserves caller traversal order for first-status selection.
	resolved []error
	// remaining bounds total work across all roots, wrappers and joins.
	remaining int
}

// visit accepts only explicit classifications or completely resolved subtrees.
func (r *errorResolver) visit(err error) bool {
	if isNilValue(err) || r.remaining == 0 {
		return false
	}
	r.remaining--
	if reflect.ValueOf(err).Comparable() {
		if _, exists := r.manifest[err]; exists {
			r.append(err)
			return true
		}
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		branches := joined.Unwrap()
		if len(branches) == 0 || len(branches) > r.remaining {
			return false
		}
		for _, branch := range branches {
			if !r.visit(branch) {
				return false
			}
		}
		return true
	}
	if classified, ok := err.(interface{ Is(error) bool }); ok {
		var match error
		for key := range r.manifest {
			if !isNilValue(key) && classified.Is(key) {
				if match != nil {
					return false
				}
				match = key
			}
		}
		if match != nil {
			r.append(match)
			return true
		}
	}
	if wrapped, ok := err.(interface{ Unwrap() error }); ok {
		return r.visit(wrapped.Unwrap())
	}
	return false
}

// append emits an identity once while preserving first occurrence order.
func (r *errorResolver) append(key error) {
	if _, exists := r.seen[key]; exists {
		return
	}
	r.seen[key] = struct{}{}
	r.resolved = append(r.resolved, key)
}

// isNilValue also rejects typed nils before calling user-provided methods.
func isNilValue(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// samePrototype detects a pointer factory returning its shared prototype. Other
// hidden aliasing remains the custom implementation's responsibility.
func samePrototype(value, prototype interface{}) bool {
	return !isNilValue(value) && reflect.ValueOf(value).Kind() == reflect.Pointer && value == prototype
}
