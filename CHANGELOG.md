# Changelog

## Unreleased

- Pin Go 1.26.10 for development, raise the module minimum from Go 1.20 to 1.26.10, and update `testify` 1.7.0 to 1.12.1. Replace its old YAML dependency with `go.yaml.in/yaml/v3` 3.0.5 and make wrapped-error tests insensitive to error ordering.

## v2.0.0 — 2025-08-03

- Changed the module path to `github.com/ooaklee/reply/v2` and changed error manifests to use errors with `errors.Is` matching.

## v1.1.0 — 2025-07-13

- Added support for wrapped and joined errors in `NewHTTPErrorResponse` and updated the Go directive to 1.20.

## v1.0.0 — 2021-09-11

- Changed the response envelope to use `meta`, `errors` and `data`, added multi-error responses and custom error transfer objects, and generalised token fields.

## v0.2.0 — 2021-09-04

- Added response helper functions and fixed error manifest merging.

## v0.1.1 — 2021-08-31

- Added licensing information and corrected documentation.

## v0.1.0 — 2021-08-29

- Introduced the library and its response and error-manifest handling.

See [release notes](release-notes.md) for prerelease details. Dates above follow the local tag commits.
