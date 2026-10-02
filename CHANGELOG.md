# Changelog

## Unreleased

### Added

- Structural `Replier.ResolveErrors` for wrapped and joined errors, using the
  premerged manifest and a bounded traversal.
- Optional request context and a privacy-safe unmapped-response observer.

### Fixed

- Update the YAML dependency used by the assertion test helpers to `v3.0.1`,
  addressing CVE-2022-28948. Reply's production package does not import YAML.
- Correctly mapped wrapped errors no longer miss their manifest entry.
- Response state is request-local, allowing shared repliers with independent,
  concurrency-safe default or custom transfer objects.
- Invalid response/prototype inputs return errors; native writer/encoder causes
  remain discoverable through `errors.Is` and `errors.As`.

### Changed

- Joined and explicit errors now have deterministic input order and unique
  identities. Incomplete trees return a generic 500 without hiding unknown
  branches. Raw fallback logging is removed; use the opt-in observer instead.
- Custom transfer factories must return independent non-nil state, not their
  prototype pointer. Nil options and attributes are ignored.

See [the integration and migration guide](UPGRADING.md) for compatibility details,
commit pinning, context propagation and custom-object ownership requirements.
