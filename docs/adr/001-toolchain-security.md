# ADR 001: Maintain the Go toolchain and dependencies

## Context

The library's Go toolchain and test dependencies need maintenance as supported toolchains and security advisories change. The v2 module path is part of the public API and a major-version migration would affect consumers.

## Decision

Pin Go 1.26.10 in `.tool-versions`, raise the module minimum from Go 1.20 to 1.26.10, update `testify` 1.7.0 to 1.12.1, keep the Go module and its dependency graph reproducible, and update compatible dependencies with their checksums. Preserve the v2 module path. Check dependencies and toolchains against upstream sources and advisories before future updates.

## Consequences

Maintainers must install the pinned Go version and run tests, module verification and vulnerability checks when changing the toolchain or dependencies. A clean automated scan is a useful signal, not a guarantee of safety. Consumers should review the module's minimum Go version when upgrading.
